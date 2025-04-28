package es

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

// SearchResult represents the structure of the Elasticsearch search response
type SearchResult struct {
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []struct {
			Index     string          `json:"_index"`
			ID        string          `json:"_id"`
			Score     float64         `json:"_score"`
			Source    json.RawMessage `json:"_source"` // Use RawMessage to delay parsing
			Highlight json.RawMessage `json:"highlight,omitempty"`
		} `json:"hits"`
	} `json:"hits"`
}

// SearchKeeps performs a search query against the specified index alias using the provided client.
func SearchKeeps(ctx context.Context, client *elasticsearch.Client, indexAlias string, query string) (SearchResult, error) {
	nilResult := SearchResult{}

	if client == nil {
		return nilResult, fmt.Errorf("elasticsearch client is not initialized")
	}
	if indexAlias == "" {
		return nilResult, fmt.Errorf("elasticsearch index alias is not provided")
	}

	buf := BuildBody(SearchParams{
		Fields:  []string{"title", "summary", "content"},
		Keyword: query,
		Index:   indexAlias,
	})

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nilResult, fmt.Errorf("error getting response: %w", err)
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
			return nilResult, fmt.Errorf("error parsing the response body: %w", err)
		} else {
			// Print the error response body for debugging.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return nilResult, fmt.Errorf("elasticsearch search error: [%s] %s", res.Status(), e["error"].(map[string]interface{})["reason"])
		}
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nilResult, fmt.Errorf("error parsing the response body: %w", err)
	}

	log.Printf(
		"[%s] %d hits; total: %d",
		res.Status(),
		len(r.Hits.Hits),
		r.Hits.Total.Value,
	)

	return r, nil
}

// SearchMoments performs a search query against the specified moments index using the provided client.
func SearchMoments(ctx context.Context, client *elasticsearch.Client, indexAlias string, query string) (SearchResult, error) {
	nilResult := SearchResult{}

	if client == nil {
		return nilResult, fmt.Errorf("elasticsearch client is not initialized")
	}
	if indexAlias == "" {
		return nilResult, fmt.Errorf("elasticsearch index alias is not provided")
	}

	buf := BuildBody(SearchParams{
		Fields:  []string{"content", "images.description"},
		Keyword: query,
		Index:   indexAlias,
	})

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nilResult, fmt.Errorf("error getting response: %w", err)
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
			return nilResult, fmt.Errorf("error parsing the response body: %w", err)
		} else {
			// Print the error response body for debugging.
			log.Printf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
			return nilResult, fmt.Errorf("elasticsearch search error: [%s] %s", res.Status(), e["error"].(map[string]interface{})["reason"])
		}
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nilResult, fmt.Errorf("error parsing the response body: %w", err)
	}

	log.Printf(
		"[%s] %d hits; total: %d",
		res.Status(),
		len(r.Hits.Hits),
		r.Hits.Total.Value,
	)

	return r, nil
}
