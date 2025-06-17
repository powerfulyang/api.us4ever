package server

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"api.us4ever/internal/config"
	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"api.us4ever/internal/logger"
	"api.us4ever/internal/metrics"
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
	server.triggerInitialIndexing()

	// Start metrics collection
	metrics.StartMetricsCollection()

	return server
}

// triggerInitialIndexing 在后台触发初始索引过程
func (s *FiberServer) triggerInitialIndexing() {
	if s.EsClient == nil {
		esLogger.Info("skipping initial Elasticsearch indexing because client is not available")
		return
	}

	// 为keeps创建索引
	go func() {
		// check index already exist, ignore create if exists
		ctx := context.Background()
		_, err := s.EsClient.Indices.Exists([]string{s.KeepEsIndexAlias}, s.EsClient.Indices.Exists.WithContext(ctx))
		if err != nil {
			esLogger.Error("failed to check index existence",
				zap.Error(err),
			)
		} else {
			esLogger.Info(fmt.Sprintf("%s already exists, skipping initial indexing", s.KeepEsIndexAlias))
			return
		}

		// Create a background context for the initial indexing
		// Use context.Background() as this is not tied to a specific request
		esLogger.Info("starting initial Elasticsearch indexing for keeps")
		if err := es.IndexKeeps(ctx, s.EsClient, s.DbClient, s.KeepEsIndexAlias); err != nil {
			esLogger.Error("initial Elasticsearch indexing for keeps failed",
				zap.Error(err),
			)
		} else {
			esLogger.Info("initial Elasticsearch indexing for keeps completed successfully")
		}
	}()

	// 为moments创建索引
	go func() {
		ctx := context.Background()
		_, err := s.EsClient.Indices.Exists([]string{s.MomentEsIndexAlias}, s.EsClient.Indices.Exists.WithContext(ctx))
		if err != nil {
			esLogger.Error("failed to check index existence",
				zap.Error(err),
			)
		} else {
			esLogger.Info(fmt.Sprintf("%s already exists, skipping initial indexing", s.MomentEsIndexAlias))
			return
		}

		esLogger.Info("starting initial Elasticsearch indexing for moments")
		if err := es.IndexMoments(ctx, s.EsClient, s.DbClient, s.MomentEsIndexAlias); err != nil {
			esLogger.Error("initial Elasticsearch indexing for moments failed",
				zap.Error(err),
			)
		} else {
			esLogger.Info("initial Elasticsearch indexing for moments completed successfully")
		}
	}()
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
