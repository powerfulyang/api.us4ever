package database

import (
	"context"
	"fmt"
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

	// GetAllMoments retrieves all Moment entities from the database.
	GetAllMoments(ctx context.Context) ([]*ent.Moment, error)

	// Close closes the database connection
	Close() error
}

// Database implements the Service interface
type Database struct {
	client *ent.Client
}

// New creates a new database service with improved error handling
func New() (Service, error) {
	// Load database configuration
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	client, err := ent.Open("postgres", dbConfig.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		client: client,
	}, nil
}

// Health checks if the database connection is healthy
func (db *Database) Health(ctx context.Context) error {
	if db.client == nil {
		return fmt.Errorf("database client is nil")
	}

	// Create a timeout context for the health check
	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Test database connectivity by performing a simple query
	// This is more reliable than trying to access the underlying sql.DB
	_, err := db.client.User.Query().Count(healthCtx)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Client returns the ent client
func (db *Database) Client() *ent.Client {
	return db.client
}

// GetAllKeeps retrieves all Keep entities.
func (db *Database) GetAllKeeps(ctx context.Context) ([]*ent.Keep, error) {
	keeps, err := db.client.Keep.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting all keeps: %w", err)
	}
	return keeps, nil
}

// GetAllMoments retrieves all Moment entities.
func (db *Database) GetAllMoments(ctx context.Context) ([]*ent.Moment, error) {
	moments, err := db.client.Moment.Query().
		WithMomentImages(func(q *ent.MomentImageQuery) {
			q.WithImage() // Eager load images
		}).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed getting all moments: %w", err)
	}
	return moments, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.client == nil {
		return nil // Already closed or never initialized
	}

	if err := db.client.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
