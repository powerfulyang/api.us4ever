package database

import (
	"api.us4ever/internal/ent/keep"
	"context"
	"fmt"
	"log"
	"time"

	"api.us4ever/internal/config"
	"api.us4ever/internal/ent"

	_ "github.com/lib/pq"
)

// Service defines the interface for database operations
type Service interface {
	// Health checks if the database connection is healthy
	Health(ctx context.Context) error

	// Client returns the ent client
	Client() *ent.Client

	// GetAllKeeps retrieves all Keep entities from the database.
	GetAllKeeps(ctx context.Context) ([]*ent.Keep, error)

	// Close closes the database connection
	Close() error
}

// Database implements the Service interface
type Database struct {
	client *ent.Client
	config *config.DBConfig
}

// New creates a new database service
func New() (Service, error) {
	// Load database configuration
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %v", err)
	}

	// Create ent client
	client, err := ent.Open("postgres", dbConfig.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &Database{
		client: client,
		config: dbConfig,
	}, nil
}

// Health checks if the database connection is healthy
func (db *Database) Health(ctx context.Context) error {
	// Create a timeout context for the health check
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try to query something simple to verify connection is working
	// Using the User entity as an example - just count the records
	_, err := db.client.User.Query().Count(ctx)
	if err != nil {
		return fmt.Errorf("database health check failed: %v", err)
	}

	return nil
}

// Client returns the ent client
func (db *Database) Client() *ent.Client {
	return db.client
}

// GetAllKeeps retrieves all Keep entities.
func (db *Database) GetAllKeeps(ctx context.Context) ([]*ent.Keep, error) {
	keeps, err := db.client.Keep.Query().Where(keep.IsPublic(true)).All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting all keeps: %w", err)
	}
	return keeps, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.client.Close()
}

// RefreshConnection recreates the database connection based on updated configuration
func RefreshConnection(db *Database) (Service, error) {
	// Close the existing connection
	if err := db.Close(); err != nil {
		log.Printf("Warning: error closing previous database connection: %v", err)
	}

	// Create a new connection
	return New()
}
