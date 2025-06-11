package tools

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"api.us4ever/internal/logger"
	"github.com/google/uuid"
)

var (
	toolsLogger *logger.Logger
)

func init() {
	var err error
	toolsLogger, err = logger.New("tools")
	if err != nil {
		panic("failed to initialize tools logger: " + err.Error())
	}
}

// ImportMomentsFromCSV imports data from CSV file to moment table
func ImportMomentsFromCSV(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path is required")
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Initialize database service
	db, err := database.New()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			toolsLogger.Warn("failed to close database connection", logger.Fields{
				"error": closeErr.Error(),
			})
		}
	}()

	// Open CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			toolsLogger.Warn("failed to close CSV file", logger.Fields{
				"error": closeErr.Error(),
			})
		}
	}()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return fmt.Errorf("CSV file is empty")
		}
		return fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Find content field index
	contentIndex := -1
	for i, header := range headers {
		if header == "content" {
			contentIndex = i
			break
		}
	}
	if contentIndex == -1 {
		return fmt.Errorf("CSV file missing required 'content' field")
	}

	// Get user ID for ownership
	userID, err := db.Client().User.Query().FirstID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	// Process CSV records
	var processedCount int
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break // End of file
			}
			toolsLogger.Warn("error reading CSV record", logger.Fields{
				"error": err.Error(),
			})
			continue
		}

		// Validate record length
		if contentIndex >= len(record) {
			toolsLogger.Warn("skipping record with insufficient columns", logger.Fields{
				"expected_columns": contentIndex + 1,
				"actual_columns":   len(record),
			})
			continue
		}

		content := record[contentIndex]
		if content == "" {
			continue
		}

		// Generate content vector
		contentVector, err := es.Embed(ctx, content)
		if err != nil {
			toolsLogger.Warn("failed to generate embedding for content", logger.Fields{
				"error":   err.Error(),
				"content": content[:min(50, len(content))], // Log first 50 chars
			})
			// Continue without vector if embedding fails
			contentVector = nil
		}

		var vectorJSON []byte
		if contentVector != nil {
			vectorJSON, err = json.Marshal(contentVector)
			if err != nil {
				toolsLogger.Warn("failed to marshal content vector", logger.Fields{
					"error": err.Error(),
				})
				vectorJSON = nil
			}
		}

		// Create new moment
		now := time.Now()
		momentID := uuid.New().String()

		_, err = db.Client().Moment.Create().
			SetID(momentID).
			SetContent(content).
			SetContentVector(vectorJSON).
			SetCategory("default").
			SetIsPublic(true).
			SetLikes(0).
			SetViews(0).
			SetTags(json.RawMessage("[]")).
			SetExtraData(json.RawMessage("{}")).
			SetOwnerId(userID).
			SetCreatedAt(now).
			SetUpdatedAt(now).
			Save(ctx)

		if err != nil {
			return fmt.Errorf("failed to save moment: %w", err)
		}

		processedCount++
		if processedCount%100 == 0 {
			toolsLogger.Info("processing progress", logger.Fields{
				"processed_count": processedCount,
			})
		}
	}

	toolsLogger.Info("CSV import completed successfully", logger.Fields{
		"total_processed": processedCount,
		"file_path":       filePath,
	})
	return nil
}
