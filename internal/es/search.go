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

	tv, _ := Embed(ctx, query)
	sv, _ := Embed(ctx, query)
	cv, _ := Embed(ctx, query)

	body := map[string]any{
		"knn": []any{
			map[string]any{"field": "title_vector", "query_vector": tv, "k": 20, "num_candidates": 60, "boost": 3},
			map[string]any{"field": "summary_vector", "query_vector": sv, "k": 20, "num_candidates": 60, "boost": 2},
			map[string]any{"field": "content_vector", "query_vector": cv, "k": 30, "num_candidates": 100, "boost": 1},
		},
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  query,
				"fields": []string{"title^3", "summary^2", "content"},
				"type":   "best_fields",
				"slop":   1,
			},
		},
		"highlight": map[string]any{
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
			"fields": map[string]any{
				"title": map[string]any{
					"number_of_fragments": 0,
				},
				"summary": map[string]any{
					"number_of_fragments": 0,
				},
				"content": map[string]any{
					"type": "unified",
					"highlight_query": map[string]any{
						"match_phrase": map[string]any{
							"content": map[string]any{
								"query": query,
								"slop":  2,
							},
						},
					},
					"number_of_fragments": 0,
				},
			},
		},
		"size": 10,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return SearchResult{}, err
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(&buf),
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

	cv, _ := Embed(ctx, query)

	body := map[string]any{
		"knn": []any{
			map[string]any{"field": "content_vector", "query_vector": cv, "k": 30, "num_candidates": 100, "boost": 1},
		},
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  query,
				"fields": []string{"content", "images.description"},
				"type":   "best_fields",
				"slop":   1,
			},
		},
		"highlight": map[string]any{
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
			"fields": map[string]any{
				"content": map[string]any{
					"type": "unified",
					"highlight_query": map[string]any{
						"match_phrase": map[string]any{
							"content": map[string]any{
								"query": query,
								"slop":  2,
							},
						},
					},
					"number_of_fragments": 0,
				},
			},
		},
		"size": 10,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return SearchResult{}, err
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(&buf),
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
