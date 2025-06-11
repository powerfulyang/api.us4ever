package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"api.us4ever/internal/logger"
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	searchLogger *logger.Logger
)

func init() {
	var err error
	searchLogger, err = logger.New("search")
	if err != nil {
		panic("failed to initialize search logger: " + err.Error())
	}
}

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

	tv, err := Embed(ctx, query)
	if err != nil {
		return nilResult, fmt.Errorf("embedding error: %w", err)
	}
	sv, err := Embed(ctx, query)
	if err != nil {
		return nilResult, fmt.Errorf("embedding error: %w", err)
	}
	cv, err := Embed(ctx, query)
	if err != nil {
		return nilResult, fmt.Errorf("embedding error: %w", err)
	}

	body := map[string]any{
		// 语义召回（保持不变）
		"knn": []any{
			map[string]any{
				"field":          "title_vector",
				"query_vector":   tv,
				"k":              20,
				"num_candidates": 60,
				"boost":          7,
			},
			map[string]any{
				"field":          "summary_vector",
				"query_vector":   sv,
				"k":              20,
				"num_candidates": 60,
				"boost":          6,
			},
			map[string]any{
				"field":          "content_vector",
				"query_vector":   cv,
				"k":              30,
				"num_candidates": 100,
				"boost":          5,
			},
		},
		"_source": map[string]any{
			"excludes": []string{"title_vector", "summary_vector", "content_vector"},
		},
		// 关键词 + 短语两路并行
		"query": map[string]any{
			"bool": map[string]any{
				"should": []any{
					// ① 两个词都得出现
					map[string]any{
						"multi_match": map[string]any{
							"query":    query,
							"fields":   []string{"title^3", "summary^2", "content"},
							"type":     "best_fields",
							"operator": "and", // 两词必须都在
							"boost":    3,     // 关键词 boost
						},
					},
					// ② 强力短语 boost
					map[string]any{
						"multi_match": map[string]any{
							"query":  query,
							"fields": []string{"title^3", "summary^2", "content"},
							"type":   "phrase",
							"slop":   2, // 允许错位
							"boost":  5, // 短语 boost
						},
					},
				},
				// should 至少命中一条即可
				"minimum_should_match": 1,
			},
		},
		// ③ 高亮：跟短语完全一致
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
	err = json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return SearchResult{}, err
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(&buf),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nilResult, fmt.Errorf("error getting response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			searchLogger.Error("error closing response body", logger.Fields{
				"error": err.Error(),
			})
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nilResult, fmt.Errorf("error parsing the response body: %w", err)
		} else {
			// Print the error response body for debugging.
			searchLogger.Error("elasticsearch search error", logger.Fields{
				"status": res.Status(),
				"type":   e["error"].(map[string]interface{})["type"],
				"reason": e["error"].(map[string]interface{})["reason"],
			})
			return nilResult, fmt.Errorf("elasticsearch search error: [%s] %s", res.Status(), e["error"].(map[string]interface{})["reason"])
		}
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nilResult, fmt.Errorf("error parsing the response body: %w", err)
	}

	searchLogger.Info("search completed", logger.Fields{
		"status":     res.Status(),
		"hits_count": len(r.Hits.Hits),
		"total":      r.Hits.Total.Value,
	})

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

	cv, err := Embed(ctx, query)
	if err != nil {
		return nilResult, fmt.Errorf("embedding error: %w", err)
	}

	body := map[string]any{
		// 语义召回（保持不变）
		"knn": []any{
			map[string]any{
				"field":          "content_vector",
				"query_vector":   cv,
				"k":              30,
				"num_candidates": 100,
				"boost":          5, // 语义召回权重
			},
		},
		"_source": map[string]any{
			"excludes": []string{"content_vector"},
		},
		// 关键词 + 短语两路并行
		"query": map[string]any{
			"bool": map[string]any{
				"should": []any{
					// ① 两个词都得出现
					map[string]any{
						"multi_match": map[string]any{
							"query":    query,
							"fields":   []string{"content", "images.description"},
							"type":     "best_fields",
							"operator": "and", // 两词必须都在
							"boost":    3,     // 关键词 boost
						},
					},
					// ② 强力短语 boost
					map[string]any{
						"multi_match": map[string]any{
							"query":  query,
							"fields": []string{"content", "images.description"},
							"type":   "phrase",
							"slop":   2, // 允许错位
							"boost":  5, // 短语 boost
						},
					},
				},
				// should 至少命中一条即可
				"minimum_should_match": 1,
			},
		},
		// ③ 高亮：跟短语完全一致
		"highlight": map[string]any{
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
			"fields": map[string]any{
				"content": map[string]any{
					"highlight_query": map[string]any{
						"match_phrase": map[string]any{
							"content": map[string]any{
								"query": query,
								"slop":  2, // 允许错位
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
	err = json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nilResult, err
	}

	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(indexAlias), // Use the provided index alias
		client.Search.WithBody(&buf),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nilResult, fmt.Errorf("error getting response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			searchLogger.Error("error closing response body", logger.Fields{
				"error": err.Error(),
			})
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nilResult, fmt.Errorf("error parsing the response body: %w", err)
		} else {
			// Print the error response body for debugging.
			searchLogger.Error("elasticsearch search error for moments", logger.Fields{
				"status": res.Status(),
				"type":   e["error"].(map[string]interface{})["type"],
				"reason": e["error"].(map[string]interface{})["reason"],
			})
			return nilResult, fmt.Errorf("elasticsearch search error: [%s] %s", res.Status(), e["error"].(map[string]interface{})["reason"])
		}
	}

	var r SearchResult
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nilResult, fmt.Errorf("error parsing the response body: %w", err)
	}

	searchLogger.Info("moments search completed", logger.Fields{
		"status":     res.Status(),
		"hits_count": len(r.Hits.Hits),
		"total":      r.Hits.Total.Value,
	})

	return r, nil
}
