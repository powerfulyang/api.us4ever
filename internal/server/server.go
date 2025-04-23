package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"api.us4ever/internal/config"
	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App

	db                 database.Service
	esClient           *elasticsearch.Client
	keepEsIndexAlias   string
	momentEsIndexAlias string
	cfg                *config.AppConfig
}

func New() *FiberServer {
	// 获取配置
	appConfig := config.GetAppConfig()

	// 初始化数据库服务
	db, err := database.New()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化 Elasticsearch 客户端
	esClient, err := es.NewClient(appConfig.ES)
	if err != nil {
		// Log the error but allow the server to start without ES if needed
		log.Printf("Failed to initialize Elasticsearch client: %v. Search functionality might be unavailable.", err)
		esClient = nil // Ensure esClient is nil if initialization fails
	}

	// Define a unique index alias (e.g., appname-keeps)
	// Ensure AppName is sanitized if it can contain special characters
	keepIndexAlias := fmt.Sprintf("%s-keeps", strings.ToLower(strings.ReplaceAll(appConfig.AppName, " ", "-")))
	momentIndexAlias := fmt.Sprintf("%s-moments", strings.ToLower(strings.ReplaceAll(appConfig.AppName, " ", "-")))

	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "api.us4ever",
			AppName:      "api.us4ever",
		}),

		db:                 db,
		esClient:           esClient,
		keepEsIndexAlias:   keepIndexAlias,
		momentEsIndexAlias: momentIndexAlias,
		cfg:                appConfig,
	}

	// 注册配置变更回调
	config.RegisterChangeCallback(server.handleConfigChange)

	// Trigger initial indexing in the background if ES client is available
	if server.esClient != nil {
		go func() {
			// Create a background context for the initial indexing
			// Use context.Background() as this is not tied to a specific request
			ctx := context.Background()
			log.Println("Starting initial Elasticsearch indexing for keeps...")
			if err := es.IndexKeeps(ctx, server.esClient, server.db, server.keepEsIndexAlias); err != nil {
				log.Printf("Initial Elasticsearch indexing for keeps failed: %v", err)
			} else {
				log.Println("Initial Elasticsearch indexing for keeps completed successfully.")
			}
		}()

		// Start a separate goroutine for indexing moments
		go func() {
			// Create a background context for the initial indexing
			ctx := context.Background()
			log.Println("Starting initial Elasticsearch indexing for moments...")
			if err := es.IndexMoments(ctx, server.esClient, server.db, server.momentEsIndexAlias); err != nil {
				log.Printf("Initial Elasticsearch indexing for moments failed: %v", err)
			} else {
				log.Println("Initial Elasticsearch indexing for moments completed successfully.")
			}
		}()
	} else {
		log.Println("Skipping initial Elasticsearch indexing because client is not available.")
	}

	return server
}

// handleConfigChange 处理配置变更
func (s *FiberServer) handleConfigChange(newConfig *config.AppConfig) {
	log.Println("配置变更，检查是否需要更新服务...")

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
		log.Println("数据库配置变更，正在更新数据库连接...")
		if err := s.RefreshDatabase(); err != nil {
			log.Printf("更新数据库连接失败: %v", err)
		} else {
			log.Println("数据库连接已更新")
		}
	} else {
		log.Println("数据库配置未变更，跳过数据库连接刷新。")
	}

	// Only refresh the ES client if the ES config actually changed
	if esConfigChanged {
		log.Println("Elasticsearch 配置变更，正在更新 ES 客户端...")
		if err := s.RefreshESClient(); err != nil {
			log.Printf("更新 Elasticsearch 客户端失败: %v", err)
		} else {
			log.Println("Elasticsearch 客户端已更新")
		}
	} else {
		log.Println("Elasticsearch 配置未变更，跳过 ES 客户端刷新。")
	}
}

// RefreshDatabase 重新创建数据库连接
func (s *FiberServer) RefreshDatabase() error {
	// 不使用类型断言，直接创建新连接
	newDb, err := database.New()
	if err != nil {
		return err
	}

	// 如果有旧连接，尝试关闭
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Printf("Warning: error closing previous database connection: %v", err)
		}
	}

	// 更新连接
	s.db = newDb
	return nil
}

// RefreshESClient 重新创建 Elasticsearch 客户端连接
func (s *FiberServer) RefreshESClient() error {
	newESClient, err := es.NewClient(s.cfg.ES)
	if err != nil {
		// Log the error but keep the old client if creation fails?
		// Or set to nil to indicate failure?
		// Setting to nil for now to indicate unavailability.
		s.esClient = nil
		return fmt.Errorf("failed to create new Elasticsearch client: %w", err)
	}

	// No explicit close needed for the standard http transport used by default.
	// If a custom transport needing cleanup is used later, add close logic here.
	s.esClient = newESClient
	return nil
}

// GetDB 返回数据库服务实例
func (s *FiberServer) GetDB() database.Service {
	return s.db
}

// ReindexKeepsHandler triggers the re-indexing process for keeps.
func (s *FiberServer) ReindexKeepsHandler(c *fiber.Ctx) error {
	log.Println("Received request to re-index keeps.")
	if s.esClient == nil {
		log.Println("ReindexKeepsHandler: Elasticsearch client is not available.")
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Elasticsearch service is not available to perform re-indexing",
		})
	}

	// Run indexing in a goroutine to avoid blocking the request
	go func() {
		// Use a background context detached from the HTTP request
		ctx := context.Background()
		log.Println("Starting background re-indexing process...")
		if err := es.IndexKeeps(ctx, s.esClient, s.db, s.keepEsIndexAlias); err != nil {
			log.Printf("Background re-indexing failed: %v", err)
			// Consider adding monitoring/alerting here
		} else {
			log.Println("Background re-indexing completed successfully.")
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process started in the background.",
	})
}

// SearchKeepsHandler handles requests to search keeps in Elasticsearch
func (s *FiberServer) SearchKeepsHandler(c *fiber.Ctx) error {
	// Get the search query from the query parameter 'q'
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing search query parameter 'q'",
		})
	}

	// Check if the ES client is available
	if s.esClient == nil {
		log.Printf("SearchKeepsHandler: Elasticsearch client is not available.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search service is temporarily unavailable",
		})
	}

	// Perform the search using the es package, passing the server's client and alias
	// The function now returns []es.KeepSearchResult
	keeps, err := es.SearchKeeps(c.Context(), s.esClient, s.keepEsIndexAlias, query)
	if err != nil {
		log.Printf("Error searching keeps in Elasticsearch: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search keeps",
		})
	}

	return c.JSON(keeps) // Return the results (including score)
}

// SearchMomentsHandler handles requests to search moments in Elasticsearch
func (s *FiberServer) SearchMomentsHandler(c *fiber.Ctx) error {
	// Get the search query from the query parameter 'q'
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing search query parameter 'q'",
		})
	}

	// Check if the ES client is available
	if s.esClient == nil {
		log.Printf("SearchMomentsHandler: Elasticsearch client is not available.")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search service is temporarily unavailable",
		})
	}

	// Perform the search using the es package, passing the server's client and alias
	moments, err := es.SearchMoments(c.Context(), s.esClient, s.momentEsIndexAlias, query)
	if err != nil {
		log.Printf("Error searching moments in Elasticsearch: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search moments",
		})
	}

	return c.JSON(moments) // Return the results (including score)
}

// ReindexMomentsHandler handles requests to re-index moments in Elasticsearch
func (s *FiberServer) ReindexMomentsHandler(c *fiber.Ctx) error {
	log.Println("Received request to re-index moments.")
	if s.esClient == nil {
		log.Println("ReindexMomentsHandler: Elasticsearch client is not available.")
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Elasticsearch service is not available to perform re-indexing",
		})
	}

	// Run indexing in a goroutine to avoid blocking the request
	go func() {
		// Use a background context detached from the HTTP request
		ctx := context.Background()
		log.Println("Starting background re-indexing process for moments...")
		if err := es.IndexMoments(ctx, s.esClient, s.db, s.momentEsIndexAlias); err != nil {
			log.Printf("Background re-indexing for moments failed: %v", err)
			// Consider adding monitoring/alerting here
		} else {
			log.Println("Background re-indexing for moments completed successfully.")
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process for moments started in the background.",
	})
}
