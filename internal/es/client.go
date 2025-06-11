package es

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	esClientLogger *logger.Logger
)

func init() {
	var err error
	esClientLogger, err = logger.New("es-client")
	if err != nil {
		panic("failed to initialize es-client logger: " + err.Error())
	}
}

// NewClient creates and returns a new Elasticsearch client based on the provided configuration
func NewClient(cfg config.ESConfig) (*elasticsearch.Client, error) {
	// Validate configuration
	if len(cfg.Addresses) == 0 {
		return nil, fmt.Errorf("elasticsearch addresses are required")
	}

	// Prepare the Elasticsearch configuration
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: &http.Transport{
			// Increase connection pool for batch operations
			MaxIdleConnsPerHost:   10,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	// Create the Elasticsearch client
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %w", err)
	}

	// Test the connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch connection test failed: %s", res.Status())
	}

	esClientLogger.Info("Elasticsearch client created and connection verified successfully")
	return client, nil
}
