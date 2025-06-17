package routes

import (
	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"api.us4ever/internal/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

var esLogger *logger.Logger

func init() {
	var err error
	esLogger, err = logger.New("elasticsearch")
	if err != nil {
		panic("failed to initialize elasticsearch logger: " + err.Error())
	}
}

type SearchRoutes struct {
	app                *fiber.App
	esClient           *elasticsearch.Client
	dbClient           database.Service
	keepEsIndexAlias   string
	momentEsIndexAlias string
}

func NewSearchRoutes(app *fiber.App, esClient *elasticsearch.Client, dbClient database.Service, keepEsIndexAlias string, momentEsIndexAlias string) *SearchRoutes {
	return &SearchRoutes{
		app:                app,
		esClient:           esClient,
		dbClient:           dbClient,
		keepEsIndexAlias:   keepEsIndexAlias,
		momentEsIndexAlias: momentEsIndexAlias,
	}
}

func (r *SearchRoutes) Register() {
	internal := r.app.Group("/internal")
	searchGroup := internal.Group("/search")

	// 新的搜索路由
	searchGroup.Get("/keeps", r.searchKeepsHandler)
	searchGroup.Get("/moments", r.searchMomentsHandler)

	// 保持向后兼容的旧路由
	internal.Get("/keeps/search", r.searchKeepsHandler)
	internal.Get("/moments/search", r.searchMomentsHandler)
}

// searchKeepsHandler handles requests to search keeps in Elasticsearch
func (r *SearchRoutes) searchKeepsHandler(c fiber.Ctx) error {

	// Get the search query from the query parameter 'q'
	query := c.Query("q")

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

	// Check if the ES client is available
	if r.esClient == nil {
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

	// Perform the search using the es package, passing the client and alias
	keeps, err := es.SearchKeeps(c.Context(), r.esClient, r.keepEsIndexAlias, query)
	if err != nil {
		esLogger.Error("error searching keeps in Elasticsearch",
			zap.Error(err),
			zap.String("query", query),
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
	esLogger.Info("search keeps completed",
		zap.String("query", query),
		zap.Int("results", len(keeps.Hits.Hits)),
		zap.Int("total", keeps.Hits.Total.Value),
	)

	return c.JSON(keeps)
}

// searchMomentsHandler handles requests to search moments in Elasticsearch
func (r *SearchRoutes) searchMomentsHandler(c fiber.Ctx) error {

	// Get the search query from the query parameter 'q'
	query := c.Query("q")

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

	// Check if the ES client is available
	if r.esClient == nil {
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

	// Perform the search using the es package, passing the client and alias
	moments, err := es.SearchMoments(c.Context(), r.esClient, r.momentEsIndexAlias, query)
	if err != nil {
		esLogger.Error("error searching moments in Elasticsearch",
			zap.Error(err),
			zap.String("query", query),
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
	esLogger.Info("search moments completed",
		zap.String("query", query),
		zap.Int("results", len(moments.Hits.Hits)),
		zap.Int("total", moments.Hits.Total.Value),
	)

	return c.JSON(moments)
}
