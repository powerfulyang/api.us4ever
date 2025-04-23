package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// KeepSearchResult represents a Keep document returned from Elasticsearch, including the search score.
type KeepSearchResult struct {
	ID      string  `json:"id"`    // Document ID from Elasticsearch
	Score   float64 `json:"score"` // Relevance score from Elasticsearch
	Title   string  `json:"title"`
	Summary string  `json:"summary"`
	Content string  `json:"content"`
	// Add other relevant fields from your Keep model if needed
}

// SearchResult represents the structure of the Elasticsearch search response
type SearchResult struct {
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			Index  string          `json:"_index"`
			ID     string          `json:"_id"`
			Score  float64         `json:"_score"`
			Source json.RawMessage `json:"_source"` // Use RawMessage to delay parsing
		} `json:"hits"`
	} `json:"hits"`
}

// SearchKeeps performs a search query against the specified index alias using the provided client.
func SearchKeeps(ctx context.Context, client *elasticsearch.Client, indexAlias string, query string) ([]KeepSearchResult, error) {
	if client == nil {
		return nil, fmt.Errorf("elasticsearch client is not initialized")
	}
	if indexAlias == "" {
		return nil, fmt.Errorf("elasticsearch index alias is not provided")
	}

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "summary", "content"},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %w", err)
		} else {
			// Print the error response body for debugging.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return nil, fmt.Errorf("elasticsearch search error: [%s] %s", res.Status(), e["error"].(map[string]interface{})["reason"])
		}
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	// Map the results to KeepSearchResult structs
	var keeps []KeepSearchResult
	for _, hit := range r.Hits.Hits {
		var keepData map[string]interface{} // Use a map first to get source fields
		if err := json.Unmarshal(hit.Source, &keepData); err != nil {
			log.Printf("Error unmarshaling keep source %s: %v", hit.ID, err)
			continue // Skip documents that fail to parse
		}

		// Construct the KeepSearchResult
		searchResult := KeepSearchResult{
			ID:      hit.ID,    // Use the ES document ID
			Score:   hit.Score, // Get the score from the hit
			Title:   getStringField(keepData, "title"),
			Summary: getStringField(keepData, "summary"),
			Content: getStringField(keepData, "content"),
			// Map other fields if necessary
		}
		keeps = append(keeps, searchResult)
	}

	log.Printf(
		"[%s] %d hits; total: %d",
		res.Status(),
		len(r.Hits.Hits),
		r.Hits.Total.Value,
	)

	return keeps, nil
}

// Helper function to safely get string fields from the map[string]interface{}
func getStringField(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return "" // Return empty string if key not found or not a string
}

// MomentSearchResult represents a Moment document returned from Elasticsearch, including the search score.
type MomentSearchResult struct {
	ID      string      `json:"id"`    // Document ID from Elasticsearch
	Score   float64     `json:"score"` // Relevance score from Elasticsearch
	Content string      `json:"content"`
	Images  []ImageInfo `json:"images,omitempty"`
}

// ImageInfo represents image information related to a moment
type ImageInfo struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
}

// SearchMoments performs a search query against the specified moments index using the provided client.
func SearchMoments(ctx context.Context, client *elasticsearch.Client, indexAlias string, query string) ([]MomentSearchResult, error) {
	if client == nil {
		return nil, fmt.Errorf("elasticsearch client is not initialized")
	}
	if indexAlias == "" {
		return nil, fmt.Errorf("elasticsearch index alias is not provided")
	}

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"content", "images.description"},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %w", err)
		} else {
			// Print the error response body for debugging.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return nil, fmt.Errorf("elasticsearch search error: [%s] %s", res.Status(), e["error"].(map[string]interface{})["reason"])
		}
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	// Map the results to MomentSearchResult structs
	var moments []MomentSearchResult
	for _, hit := range r.Hits.Hits {
		var momentData map[string]interface{} // Use a map first to get source fields
		if err := json.Unmarshal(hit.Source, &momentData); err != nil {
			log.Printf("Error unmarshaling moment source %s: %v", hit.ID, err)
			continue // Skip documents that fail to parse
		}

		// Process images if available
		var images []ImageInfo
		if imagesData, hasImages := momentData["images"].([]interface{}); hasImages {
			for _, imgData := range imagesData {
				if img, ok := imgData.(map[string]interface{}); ok {
					imageInfo := ImageInfo{
						ID:          getStringField(img, "id"),
						Description: getStringField(img, "description"),
					}
					images = append(images, imageInfo)
				}
			}
		}

		// Construct the MomentSearchResult
		searchResult := MomentSearchResult{
			ID:      hit.ID,
			Score:   hit.Score,
			Content: getStringField(momentData, "content"),
			Images:  images,
		}
		moments = append(moments, searchResult)
	}

	log.Printf(
		"[%s] %d hits; total: %d",
		res.Status(),
		len(r.Hits.Hits),
		r.Hits.Total.Value,
	)

	return moments, nil
}
