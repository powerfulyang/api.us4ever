package es

import (
	"fmt"
	"log"

	"api.us4ever/internal/config"
	"github.com/elastic/go-elasticsearch/v8"
)

// NewClient creates and returns a new Elasticsearch client based on the provided configuration.
func NewClient(cfg config.ESConfig) (*elasticsearch.Client, error) {
	// Prepare the Elasticsearch configuration
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	// Create the Elasticsearch client
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		log.Printf("Error creating the Elasticsearch client: %s", err)
		return nil, fmt.Errorf("error creating the Elasticsearch client: %w", err)
	}

	// Test the connection
	_, err = client.Info()
	if err != nil {
		log.Printf("Error connecting to Elasticsearch: %s", err)
		// Don't return the client if the connection test failed
		return nil, fmt.Errorf("error connecting to Elasticsearch: %w", err)
	}

	log.Println("Elasticsearch client created and connection verified successfully")
	return client, nil
}
