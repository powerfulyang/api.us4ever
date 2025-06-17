package routes

import (
	"context"
	"net/http"

	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"api.us4ever/internal/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

var reindexLogger *logger.Logger

func init() {
	var err error
	reindexLogger, err = logger.New("elasticsearch")
	if err != nil {
		panic("failed to initialize elasticsearch logger: " + err.Error())
	}
}

type ReindexRoutes struct {
	app                *fiber.App
	esClient           *elasticsearch.Client
	dbClient           database.Service
	keepEsIndexAlias   string
	momentEsIndexAlias string
}

func NewReindexRoutes(app *fiber.App, esClient *elasticsearch.Client, dbClient database.Service, keepEsIndexAlias string, momentEsIndexAlias string) *ReindexRoutes {
	return &ReindexRoutes{
		app:                app,
		esClient:           esClient,
		dbClient:           dbClient,
		keepEsIndexAlias:   keepEsIndexAlias,
		momentEsIndexAlias: momentEsIndexAlias,
	}
}

func (r *ReindexRoutes) Register() {
	internal := r.app.Group("/internal")

	// 重索引端点
	internal.Post("/keeps/reindex", r.reindexKeepsHandler)
	internal.Post("/moments/reindex", r.reindexMomentsHandler)
}

// reindexKeepsHandler triggers the re-indexing process for keeps.
func (r *ReindexRoutes) reindexKeepsHandler(c fiber.Ctx) error {
	reindexLogger.Info("received request to re-index keeps")
	if r.esClient == nil {
		reindexLogger.Warn("Elasticsearch client is not available for re-indexing")
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Elasticsearch service is not available to perform re-indexing",
		})
	}

	// Run indexing in a goroutine to avoid blocking the request
	go func() {
		// Use a background context detached from the HTTP request
		ctx := context.Background()
		reindexLogger.Info("starting background re-indexing process",
			zap.String("index_type", "keeps"),
			zap.String("index_alias", r.keepEsIndexAlias),
		)
		if err := es.IndexKeeps(ctx, r.esClient, r.dbClient, r.keepEsIndexAlias); err != nil {
			reindexLogger.Error("background re-indexing failed",
				zap.Error(err),
				zap.String("index_type", "keeps"),
			)
			// Consider adding monitoring/alerting here
		} else {
			reindexLogger.Info("background re-indexing completed successfully",
				zap.String("index_type", "keeps"),
			)
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process started in the background.",
	})
}

// reindexMomentsHandler handles requests to re-index moments in Elasticsearch
func (r *ReindexRoutes) reindexMomentsHandler(c fiber.Ctx) error {
	reindexLogger.Info("received request to re-index moments")
	if r.esClient == nil {
		reindexLogger.Warn("Elasticsearch client is not available for re-indexing")
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Elasticsearch service is not available to perform re-indexing",
		})
	}

	// Run indexing in a goroutine to avoid blocking the request
	go func() {
		// Use a background context detached from the HTTP request
		ctx := context.Background()
		reindexLogger.Info("starting background re-indexing process",
			zap.String("index_type", "moments"),
			zap.String("index_alias", r.momentEsIndexAlias),
		)
		if err := es.IndexMoments(ctx, r.esClient, r.dbClient, r.momentEsIndexAlias); err != nil {
			reindexLogger.Error("background re-indexing failed",
				zap.Error(err),
				zap.String("index_type", "moments"),
			)
			// Consider adding monitoring/alerting here
		} else {
			reindexLogger.Info("background re-indexing completed successfully",
				zap.String("index_type", "moments"),
			)
		}
	}()

	// Immediately return success, indicating the process has started
	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "Re-indexing process for moments started in the background.",
	})
}
