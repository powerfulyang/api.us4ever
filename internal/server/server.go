package server

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"api.us4ever/internal/config"
	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"api.us4ever/internal/logger"
	"api.us4ever/internal/middleware"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type FiberServer struct {
	*fiber.App

	DbClient           database.Service
	EsClient           *elasticsearch.Client
	KeepEsIndexAlias   string
	MomentEsIndexAlias string
	cfg                *config.AppConfig
	logger             *logger.Logger
}

var (
	serverLogger *logger.Logger
	esLogger     *logger.Logger
	configLogger *logger.Logger
)

func init() {
	var err error

	serverLogger, err = logger.New("server")
	if err != nil {
		panic("failed to initialize server logger: " + err.Error())
	}

	esLogger, err = logger.New("elasticsearch")
	if err != nil {
		panic("failed to initialize elasticsearch logger: " + err.Error())
	}

	configLogger, err = logger.New("config")
	if err != nil {
		panic("failed to initialize config logger: " + err.Error())
	}
}

func New() *FiberServer {
	// Load configuration
	appConfig := config.MustGetAppConfig()

	// Initialize database service
	dbClient, err := database.New()
	if err != nil {
		serverLogger.Fatal("failed to initialize database",
			zap.Error(err),
		)
	}

	// Initialize Elasticsearch client
	var esClient *elasticsearch.Client
	if len(appConfig.ES.Addresses) > 0 {
		esClient, err = es.NewClient(appConfig.ES)
		if err != nil {
			// Log the error but allow the server to start without ES if needed
			esLogger.Error("failed to initialize Elasticsearch client",
				zap.Error(err),
			)
			esClient = nil
		}
	} else {
		esLogger.Info("Elasticsearch configuration not provided, search functionality will be unavailable")
	}

	// Create index aliases with sanitized app name
	sanitizedAppName := strings.ToLower(strings.ReplaceAll(appConfig.AppName, " ", "-"))
	keepIndexAlias := fmt.Sprintf("%s-keeps", sanitizedAppName)
	momentIndexAlias := fmt.Sprintf("%s-moments", sanitizedAppName)

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: appConfig.AppName,
			AppName:      appConfig.AppName,
			ErrorHandler: middleware.NewErrorHandler(),
		}),

		DbClient:           dbClient,
		EsClient:           esClient,
		KeepEsIndexAlias:   keepIndexAlias,
		MomentEsIndexAlias: momentIndexAlias,
		cfg:                appConfig,
	}

	// Register configuration change callback
	config.RegisterChangeCallback(server.handleConfigChange)

	// Trigger initial indexing in the background if ES client is available
	if server.EsClient != nil {
		go func() {
			// check index already exist, ignore create if exists
			ctx := context.Background()
			_, err := server.EsClient.Indices.Exists([]string{server.KeepEsIndexAlias}, server.EsClient.Indices.Exists.WithContext(ctx))
			if err != nil {
				esLogger.Error("failed to check index existence",
					zap.Error(err),
				)
			} else {
				esLogger.Info(fmt.Sprintf("%s already exists, skipping initial indexing", server.KeepEsIndexAlias))
				return
			}

			// Create a background context for the initial indexing
			// Use context.Background() as this is not tied to a specific request
			esLogger.Info("starting initial Elasticsearch indexing for keeps")
			if err := es.IndexKeeps(ctx, server.EsClient, server.DbClient, server.KeepEsIndexAlias); err != nil {
				esLogger.Error("initial Elasticsearch indexing for keeps failed",
					zap.Error(err),
				)
			} else {
				esLogger.Info("initial Elasticsearch indexing for keeps completed successfully")
			}
		}()

		// Start a separate goroutine for indexing moments
		go func() {
			ctx := context.Background()
			_, err := server.EsClient.Indices.Exists([]string{server.MomentEsIndexAlias}, server.EsClient.Indices.Exists.WithContext(ctx))
			if err != nil {
				esLogger.Error("failed to check index existence",
					zap.Error(err),
				)
			} else {
				esLogger.Info(fmt.Sprintf("%s already exists, skipping initial indexing", server.MomentEsIndexAlias))
				return
			}

			esLogger.Info("starting initial Elasticsearch indexing for moments")
			if err := es.IndexMoments(ctx, server.EsClient, server.DbClient, server.MomentEsIndexAlias); err != nil {
				esLogger.Error("initial Elasticsearch indexing for moments failed",
					zap.Error(err),
				)
			} else {
				esLogger.Info("initial Elasticsearch indexing for moments completed successfully")
			}
		}()
	} else {
		esLogger.Info("skipping initial Elasticsearch indexing because client is not available")
	}

	return server
}

// handleConfigChange handles configuration changes
func (s *FiberServer) handleConfigChange(newConfig *config.AppConfig) {
	configLogger.Info("configuration changed, checking if services need updates")

	// Store old config for comparison
	oldConfig := s.cfg

	// Update the server's config regardless
	s.cfg = newConfig

	dbConfigChanged := false
	esConfigChanged := false

	if oldConfig != nil { // Ensure current config exists for comparison
		dbConfigChanged = oldConfig.Database != newConfig.Database
		esConfigChanged = !reflect.DeepEqual(oldConfig.ES, newConfig.ES)
	}

	// Only refresh the database connection if the DB config actually changed
	if dbConfigChanged {
		configLogger.Info("database configuration changed, updating database connection")
		if err := s.refreshDatabase(); err != nil {
			configLogger.Error("failed to update database connection",
				zap.Error(err),
			)
		} else {
			configLogger.Info("database connection updated successfully")
		}
	} else {
		configLogger.Debug("database configuration unchanged, skipping database connection refresh")
	}

	// Only refresh the ES client if the ES config actually changed
	if esConfigChanged {
		esLogger.Info("Elasticsearch configuration changed, updating ES client")
		if err := s.refreshESClient(); err != nil {
			esLogger.Error("failed to update Elasticsearch client",
				zap.Error(err),
			)
		} else {
			esLogger.Info("Elasticsearch client updated successfully")
		}
	} else {
		esLogger.Debug("Elasticsearch configuration unchanged, skipping ES client refresh")
	}
}

// refreshDatabase 重新创建数据库连接
func (s *FiberServer) refreshDatabase() error {
	// 不使用类型断言，直接创建新连接
	newDb, err := database.New()
	if err != nil {
		return err
	}

	// Close old connection if exists
	if s.DbClient != nil {
		if err := s.DbClient.Close(); err != nil {
			serverLogger.Warn("error closing previous database connection",
				zap.Error(err),
			)
		}
	}

	// 更新连接
	s.DbClient = newDb
	return nil
}

// refreshESClient 重新创建 Elasticsearch 客户端连接
func (s *FiberServer) refreshESClient() error {
	newESClient, err := es.NewClient(s.cfg.ES)
	if err != nil {
		// Log the error but keep the old client if creation fails?
		// Or set to nil to indicate failure?
		// Setting to nil for now to indicate unavailability.
		s.EsClient = nil
		return fmt.Errorf("failed to create new Elasticsearch client: %w", err)
	}

	// No explicit close needed for the standard http transport used by default.
	// If a custom transport needing cleanup is used later, add close logic here.
	s.EsClient = newESClient
	return nil
}

// reindexKeepsHandler triggers the re-indexing process for keeps.
func (s *FiberServer) reindexKeepsHandler(c fiber.Ctx) error {
	esLogger.Info("received request to re-index keeps")
	if s.EsClient == nil {
		esLogger.Warn("Elasticsearch client is not available for re-indexing")
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Elasticsearch service is not available to perform re-indexing",
		})
	}

	// Run indexing in a goroutine to avoid blocking the request
	go func() {
		// Use a background context detached from the HTTP request
		ctx := context.Background()
		esLogger.Info("starting background re-indexing process",
			zap.String("index_type", "keeps"),
			zap.String("index_alias", s.KeepEsIndexAlias),
		)
		if err := es.IndexKeeps(ctx, s.EsClient, s.DbClient, s.KeepEsIndexAlias); err != nil {
			esLogger.Error("background re-indexing failed",
				zap.Error(err),
				zap.String("index_type", "keeps"),
			)
			// Consider adding monitoring/alerting here
		} else {
			esLogger.Info("background re-indexing completed successfully",
				zap.String("index_type", "keeps"),
			)
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process started in the background.",
	})
}

// searchKeepsHandler handles requests to search keeps in Elasticsearch
func (s *FiberServer) searchKeepsHandler(c fiber.Ctx) error {
	start := time.Now()

	// Get the search query from the query parameter 'q'
	query := c.Query("q")
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)

	// Basic input validation
	if query == "" {
		esLogger.Warn("search request with empty query",
			zap.String("ip", c.IP()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "ValidationError",
				"message": "Missing search query parameter 'q'",
				"code":    400,
			},
		})
	}

	// Validate query length
	if len(query) > 200 {
		esLogger.Warn("search request with query too long",
			zap.String("ip", c.IP()),
			zap.Int("query_length", len(query)),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "ValidationError",
				"message": "Search query too long (max 200 characters)",
				"code":    400,
			},
		})
	}

	// Validate limit
	if limit < 1 || limit > 100 {
		limit = 10 // Set default
	}

	// Validate offset
	if offset < 0 {
		offset = 0
	}

	// Basic sanitization - remove dangerous patterns
	if containsDangerousPatterns(query) {
		esLogger.Warn("search request with potentially dangerous content",
			zap.String("ip", c.IP()),
			zap.String("query", query),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "SecurityError",
				"message": "Search query contains invalid content",
				"code":    400,
			},
		})
	}

	// Check if the ES client is available
	if s.EsClient == nil {
		esLogger.Warn("Elasticsearch client is not available for search",
			zap.String("handler", "searchKeeps"),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "ServiceError",
				"message": "Search service is temporarily unavailable",
				"code":    503,
			},
		})
	}

	// Perform the search using the es package, passing the server's client and alias
	keeps, err := es.SearchKeeps(c.Context(), s.EsClient, s.KeepEsIndexAlias, query)
	if err != nil {
		duration := time.Since(start)
		esLogger.Error("error searching keeps in Elasticsearch",
			zap.Error(err),
			zap.String("query", query),
			zap.Duration("duration", duration),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "SearchError",
				"message": "Failed to search keeps",
				"code":    500,
			},
		})
	}

	// Log successful search
	duration := time.Since(start)
	esLogger.Info("search keeps completed",
		zap.String("query", query),
		zap.Int("results", len(keeps.Hits.Hits)),
		zap.Int("total", keeps.Hits.Total.Value),
		zap.Duration("duration", duration),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return c.JSON(fiber.Map{
		"query":    query,
		"total":    keeps.Hits.Total.Value,
		"results":  keeps.Hits.Hits,
		"duration": duration.Milliseconds(),
		"limit":    limit,
		"offset":   offset,
	})
}

// searchMomentsHandler handles requests to search moments in Elasticsearch
func (s *FiberServer) searchMomentsHandler(c fiber.Ctx) error {
	start := time.Now()

	// Get the search query from the query parameter 'q'
	query := c.Query("q")
	limit := fiber.Query[int](c, "limit", 10)
	offset := fiber.Query[int](c, "offset", 0)

	// Basic input validation
	if query == "" {
		esLogger.Warn("search request with empty query",
			zap.String("ip", c.IP()),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "ValidationError",
				"message": "Missing search query parameter 'q'",
				"code":    400,
			},
		})
	}

	// Validate query length
	if len(query) > 200 {
		esLogger.Warn("search request with query too long",
			zap.String("ip", c.IP()),
			zap.Int("query_length", len(query)),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "ValidationError",
				"message": "Search query too long (max 200 characters)",
				"code":    400,
			},
		})
	}

	// Validate limit
	if limit < 1 || limit > 100 {
		limit = 10 // Set default
	}

	// Validate offset
	if offset < 0 {
		offset = 0
	}

	// Basic sanitization - remove dangerous patterns
	if containsDangerousPatterns(query) {
		esLogger.Warn("search request with potentially dangerous content",
			zap.String("ip", c.IP()),
			zap.String("query", query),
		)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "SecurityError",
				"message": "Search query contains invalid content",
				"code":    400,
			},
		})
	}

	// Check if the ES client is available
	if s.EsClient == nil {
		esLogger.Warn("Elasticsearch client is not available for search",
			zap.String("handler", "searchMoments"),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "ServiceError",
				"message": "Search service is temporarily unavailable",
				"code":    503,
			},
		})
	}

	// Perform the search using the es package, passing the server's client and alias
	moments, err := es.SearchMoments(c.Context(), s.EsClient, s.MomentEsIndexAlias, query)
	if err != nil {
		duration := time.Since(start)
		esLogger.Error("error searching moments in Elasticsearch",
			zap.Error(err),
			zap.String("query", query),
			zap.Duration("duration", duration),
		)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fiber.Map{
				"type":    "SearchError",
				"message": "Failed to search moments",
				"code":    500,
			},
		})
	}

	// Log successful search
	duration := time.Since(start)
	esLogger.Info("search moments completed",
		zap.String("query", query),
		zap.Int("results", len(moments.Hits.Hits)),
		zap.Int("total", moments.Hits.Total.Value),
		zap.Duration("duration", duration),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	return c.JSON(fiber.Map{
		"query":    query,
		"total":    moments.Hits.Total.Value,
		"results":  moments.Hits.Hits,
		"duration": duration.Milliseconds(),
		"limit":    limit,
		"offset":   offset,
	})
}

// reindexMomentsHandler handles requests to re-index moments in Elasticsearch
func (s *FiberServer) reindexMomentsHandler(c fiber.Ctx) error {
	esLogger.Info("received request to re-index moments")
	if s.EsClient == nil {
		esLogger.Warn("Elasticsearch client is not available for re-indexing")
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Elasticsearch service is not available to perform re-indexing",
		})
	}

	// Run indexing in a goroutine to avoid blocking the request
	go func() {
		// Use a background context detached from the HTTP request
		ctx := context.Background()
		esLogger.Info("starting background re-indexing process",
			zap.String("index_type", "moments"),
			zap.String("index_alias", s.MomentEsIndexAlias),
		)
		if err := es.IndexMoments(ctx, s.EsClient, s.DbClient, s.MomentEsIndexAlias); err != nil {
			esLogger.Error("background re-indexing failed",
				zap.Error(err),
				zap.String("index_type", "moments"),
			)
			// Consider adding monitoring/alerting here
		} else {
			esLogger.Info("background re-indexing completed successfully",
				zap.String("index_type", "moments"),
			)
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process for moments started in the background.",
	})
}

// containsDangerousPatterns checks if the query contains potentially dangerous patterns
func containsDangerousPatterns(query string) bool {
	// List of dangerous patterns to check for
	dangerousPatterns := []string{
		"<script",
		"javascript:",
		"onload=",
		"onerror=",
		"eval(",
		"document.",
		"window.",
		".innerHTML",
		".outerHTML",
	}

	lowerQuery := strings.ToLower(query)

	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerQuery, strings.ToLower(pattern)) {
			return true
		}
	}

	return false
}
