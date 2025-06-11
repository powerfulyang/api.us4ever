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
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
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
		serverLogger.Fatal("failed to initialize database", logger.Fields{
			"error": err.Error(),
		})
	}

	// Initialize Elasticsearch client
	var esClient *elasticsearch.Client
	if len(appConfig.ES.Addresses) > 0 {
		esClient, err = es.NewClient(appConfig.ES)
		if err != nil {
			// Log the error but allow the server to start without ES if needed
			esLogger.Error("failed to initialize Elasticsearch client", logger.Fields{
				"error": err.Error(),
			})
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
			ServerHeader:     appConfig.AppName,
			AppName:          appConfig.AppName,
			DisableKeepalive: false,
			ReadTimeout:      30 * time.Second,
			WriteTimeout:     30 * time.Second,
			IdleTimeout:      60 * time.Second,
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
			// Create a background context for the initial indexing
			// Use context.Background() as this is not tied to a specific request
			ctx := context.Background()
			esLogger.Info("starting initial Elasticsearch indexing for keeps")
			if err := es.IndexKeeps(ctx, server.EsClient, server.DbClient, server.KeepEsIndexAlias); err != nil {
				esLogger.Error("initial Elasticsearch indexing for keeps failed", logger.Fields{
					"error": err.Error(),
				})
			} else {
				esLogger.Info("initial Elasticsearch indexing for keeps completed successfully")
			}
		}()

		// Start a separate goroutine for indexing moments
		go func() {
			// Create a background context for the initial indexing
			ctx := context.Background()
			esLogger.Info("starting initial Elasticsearch indexing for moments")
			if err := es.IndexMoments(ctx, server.EsClient, server.DbClient, server.MomentEsIndexAlias); err != nil {
				esLogger.Error("initial Elasticsearch indexing for moments failed", logger.Fields{
					"error": err.Error(),
				})
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
			configLogger.Error("failed to update database connection", logger.Fields{
				"error": err.Error(),
			})
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
			esLogger.Error("failed to update Elasticsearch client", logger.Fields{
				"error": err.Error(),
			})
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
			serverLogger.Warn("error closing previous database connection", logger.Fields{
				"error": err.Error(),
			})
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
func (s *FiberServer) reindexKeepsHandler(c *fiber.Ctx) error {
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
		esLogger.Info("starting background re-indexing process", logger.Fields{
			"index_type":  "keeps",
			"index_alias": s.KeepEsIndexAlias,
		})
		if err := es.IndexKeeps(ctx, s.EsClient, s.DbClient, s.KeepEsIndexAlias); err != nil {
			esLogger.Error("background re-indexing failed", logger.Fields{
				"error":      err.Error(),
				"index_type": "keeps",
			})
			// Consider adding monitoring/alerting here
		} else {
			esLogger.Info("background re-indexing completed successfully", logger.Fields{
				"index_type": "keeps",
			})
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process started in the background.",
	})
}

// searchKeepsHandler handles requests to search keeps in Elasticsearch
func (s *FiberServer) searchKeepsHandler(c *fiber.Ctx) error {
	// Get the search query from the query parameter 'q'
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing search query parameter 'q'",
		})
	}

	// Check if the ES client is available
	if s.EsClient == nil {
		esLogger.Warn("Elasticsearch client is not available for search", logger.Fields{
			"handler": "searchKeeps",
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search service is temporarily unavailable",
		})
	}

	// Perform the search using the es package, passing the server's client and alias
	// The function now returns []es.KeepSearchResult
	keeps, err := es.SearchKeeps(c.Context(), s.EsClient, s.KeepEsIndexAlias, query)
	if err != nil {
		esLogger.Error("error searching keeps in Elasticsearch", logger.Fields{
			"error": err.Error(),
			"query": query,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search keeps",
		})
	}

	return c.JSON(keeps) // Return the results (including score)
}

// searchMomentsHandler handles requests to search moments in Elasticsearch
func (s *FiberServer) searchMomentsHandler(c *fiber.Ctx) error {
	// Get the search query from the query parameter 'q'
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing search query parameter 'q'",
		})
	}

	// Check if the ES client is available
	if s.EsClient == nil {
		esLogger.Warn("Elasticsearch client is not available for search", logger.Fields{
			"handler": "searchMoments",
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search service is temporarily unavailable",
		})
	}

	// Perform the search using the es package, passing the server's client and alias
	moments, err := es.SearchMoments(c.Context(), s.EsClient, s.MomentEsIndexAlias, query)
	if err != nil {
		esLogger.Error("error searching moments in Elasticsearch", logger.Fields{
			"error": err.Error(),
			"query": query,
		})
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search moments",
		})
	}

	return c.JSON(moments) // Return the results (including score)
}

// reindexMomentsHandler handles requests to re-index moments in Elasticsearch
func (s *FiberServer) reindexMomentsHandler(c *fiber.Ctx) error {
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
		esLogger.Info("starting background re-indexing process", logger.Fields{
			"index_type":  "moments",
			"index_alias": s.MomentEsIndexAlias,
		})
		if err := es.IndexMoments(ctx, s.EsClient, s.DbClient, s.MomentEsIndexAlias); err != nil {
			esLogger.Error("background re-indexing failed", logger.Fields{
				"error":      err.Error(),
				"index_type": "moments",
			})
			// Consider adding monitoring/alerting here
		} else {
			esLogger.Info("background re-indexing completed successfully", logger.Fields{
				"index_type": "moments",
			})
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process for moments started in the background.",
	})
}
