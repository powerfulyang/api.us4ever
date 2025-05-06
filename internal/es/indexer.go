package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/ent"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/tidwall/gjson"
)

const ( // Constants for bulk indexing
	bulkIndexAction = `{ "index" : { "_index" : "%s", "_id" : "%s" } }`
	bulkFlushBytes  = 5 * 1024 * 1024 // Flush threshold 5MB
	bulkFlushItems  = 1000            // Flush threshold 1000 items
)

// IndexKeeps fetches all Keep records from the database and indexes them into a new
// Elasticsearch index, then atomically switches the alias to point to the new index.
func IndexKeeps(ctx context.Context, client *elasticsearch.Client, dbService database.Service, aliasName string) error {
	if client == nil {
		return fmt.Errorf("elasticsearch client is not initialized")
	}
	if dbService == nil {
		return fmt.Errorf("database service is not initialized")
	}
	if aliasName == "" {
		return fmt.Errorf("index alias name is required")
	}

	log.Printf("Starting re-indexing process for alias: %s", aliasName)

	// 1. Create a new index with a timestamp
	newIndexName := fmt.Sprintf("%s_%s", aliasName, time.Now().Format("20060102150405"))
	log.Printf("Creating new index: %s", newIndexName)

	mapping := map[string]any{
		"settings": map[string]any{
			"number_of_shards":   3,
			"number_of_replicas": 0,
			"max_ngram_diff":     2,
			"analysis": map[string]any{
				"tokenizer": map[string]any{
					"cjk_ngram": map[string]any{
						"type":     "ngram",
						"min_gram": 2,
						"max_gram": 4,
					},
				},
				"analyzer": map[string]any{
					"ik_cjk": map[string]any{
						"tokenizer": "ik_max_word",
					},
					"cjk_ngram_analyzer": map[string]any{
						"tokenizer": "cjk_ngram",
						"filter":    []string{"lowercase"},
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": MergeTextFields([]string{"title", "summary", "content"}),
		},
	}

	// 向量字段统一追加
	vecFields := map[string]string{
		"title_vector":   "title_vector",
		"summary_vector": "summary_vector",
		"content_vector": "content_vector",
	}
	props := mapping["mappings"].(map[string]any)["properties"].(map[string]any)
	for name := range vecFields {
		props[name] = map[string]any{
			"type":       "dense_vector",
			"dims":       vectorDims,
			"index":      true,
			"similarity": "cosine",
		}
	}

	body, _ := json.Marshal(mapping)

	res, err := client.Indices.Create(
		newIndexName,
		client.Indices.Create.WithContext(ctx),
		client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return fmt.Errorf("cannot create index %s: %w", newIndexName, err)
	}
	if res.IsError() {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(res.Body)
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("cannot create index %s: [%s] %s", newIndexName, res.Status(), string(bodyBytes))
	}
	log.Printf("Index %s created successfully", newIndexName)
	err = res.Body.Close()
	if err != nil {
		return err
	} // Close successful response body

	// 2. Fetch data from the database
	log.Printf("Fetching keeps from database...")
	keeps, err := dbService.GetAllKeeps(ctx)
	if err != nil {
		// Consider deleting the newly created index if DB fetch fails
		_, delErr := client.Indices.Delete([]string{newIndexName}, client.Indices.Delete.WithContext(ctx))
		if delErr != nil {
			log.Printf("Failed to delete temporary index %s after DB error: %v", newIndexName, delErr)
		}
		return fmt.Errorf("failed to fetch keeps from database: %w", err)
	}
	log.Printf("Fetched %d keeps from database", len(keeps))

	// 3. Bulk index the data
	log.Printf("Starting bulk indexing into %s...", newIndexName)
	if err := bulkIndexKeeps(ctx, client, newIndexName, keeps); err != nil {
		// Consider deleting the newly created index if bulk indexing fails
		_, delErr := client.Indices.Delete([]string{newIndexName}, client.Indices.Delete.WithContext(ctx))
		if delErr != nil {
			log.Printf("Failed to delete temporary index %s after bulk index error: %v", newIndexName, delErr)
		}
		return fmt.Errorf("bulk indexing failed: %w", err)
	}
	log.Printf("Bulk indexing into %s completed successfully", newIndexName)

	// 4. Atomically update the alias
	log.Printf("Updating alias %s to point to %s", aliasName, newIndexName)
	if err := updateAlias(ctx, client, aliasName, newIndexName); err != nil {
		// If alias update fails, the new index is orphaned but searchable directly.
		// Consider manual cleanup or retry logic.
		return fmt.Errorf("failed to update alias %s: %w", aliasName, err)
	}
	log.Printf("Alias %s updated successfully", aliasName)

	// 5. Delete old indices (run in background, log errors)
	go func() {
		if err := deleteOldIndices(context.Background(), client, aliasName, newIndexName); err != nil {
			// Error is already logged within deleteOldIndices or the function returned nil on logged error
			log.Printf("Background deletion of old indices encountered an issue: %v", err) // Log any unexpected error return
		}
	}()

	log.Printf("Re-indexing process for alias %s completed successfully", aliasName)
	return nil
}

// bulkIndexKeeps performs bulk indexing of Keep documents.
func bulkIndexKeeps(ctx context.Context, client *elasticsearch.Client, indexName string, keeps []*ent.Keep) error {
	var ( // Bulk buffer and counters
		buf    bytes.Buffer
		numOps int
	)

	for _, keep := range keeps {
		// Prepare meta line (action and metadata)
		meta := fmt.Sprintf(bulkIndexAction, indexName, keep.ID)
		buf.WriteString(meta)
		buf.WriteByte('\n')

		title := keep.Title
		summary := keep.Summary
		content := keep.Content
		titleVector := keep.TitleVector
		summaryVector := keep.SummaryVector
		contentVector := keep.ContentVector

		// Prepare data line (document source)
		// Convert ent.Keep to a suitable map/struct for JSON marshalling
		// Only include fields relevant for search (title, summary, content)
		doc := map[string]interface{}{
			"title":   title,   // Assuming ent.Keep has these fields
			"summary": summary, // Adjust field names as necessary
			"content": content,
			// Add other fields if needed for search or display
			"title_vector":   titleVector,
			"summary_vector": summaryVector,
			"content_vector": contentVector,
		}
		data, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling keep %s: %v. Skipping.", keep.ID, err)
			buf.Reset() // Clear the buffer for this item
			continue
		}
		buf.Write(data)
		buf.WriteByte('\n')

		numOps++

		// Flush buffer if thresholds reached
		if buf.Len() > bulkFlushBytes || numOps >= bulkFlushItems {
			if err := flushBulkBuffer(ctx, client, &buf); err != nil {
				return err // Propagate error up
			}
			numOps = 0 // Reset counter
		}
	}

	// Flush any remaining items in the buffer
	if buf.Len() > 0 {
		if err := flushBulkBuffer(ctx, client, &buf); err != nil {
			return err
		}
	}

	// Refresh the index to make changes searchable immediately
	_, err := client.Indices.Refresh(client.Indices.Refresh.WithContext(ctx), client.Indices.Refresh.WithIndex(indexName))
	if err != nil {
		log.Printf("Warning: Failed to refresh index %s after bulk indexing: %v", indexName, err)
		// Don't fail the whole process, but log the warning
	}

	return nil
}

// flushBulkBuffer sends the bulk request to Elasticsearch.
func flushBulkBuffer(ctx context.Context, client *elasticsearch.Client, buf *bytes.Buffer) error {
	res, err := client.Bulk(bytes.NewReader(buf.Bytes()), client.Bulk.WithContext(ctx))
	buf.Reset() // Reset buffer regardless of outcome
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(res.Body)

	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read bulk response body: %w", readErr)
	}

	if res.IsError() {
		return fmt.Errorf("bulk request returned error: [%s] %s", res.Status(), string(bodyBytes))
	}

	// Check for item-level errors in the bulk response
	jsonResponse := string(bodyBytes)
	if gjson.Get(jsonResponse, "errors").Bool() {
		var errorCount int
		gjson.Get(jsonResponse, "items").ForEach(func(key, value gjson.Result) bool {
			if value.Get("index.error").Exists() {
				errorCount++
				log.Printf("Item error: Index=%s, ID=%s, Status=%d, Type=%s, Reason=%s",
					value.Get("index._index").String(),
					value.Get("index._id").String(),
					value.Get("index.status").Int(),
					value.Get("index.error.type").String(),
					value.Get("index.error.reason").String(),
				)
			}
			return true // continue iterating
		})
		return fmt.Errorf("bulk request completed with %d item errors (see logs for details)", errorCount)
	}

	log.Printf("Bulk buffer flushed successfully.")
	return nil
}

// updateAlias atomically switches the alias to point to the new index.
func updateAlias(ctx context.Context, client *elasticsearch.Client, aliasName, newIndexName string) error {
	body := map[string][]map[string]map[string]interface{}{
		"actions": {
			{
				"remove": {
					"alias": aliasName,
					"index": aliasName + "_*", // Remove from all indices matching the pattern
				},
			},
			{
				"add": {
					"alias": aliasName,
					"index": newIndexName,
				},
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal alias update body: %w", err)
	}

	res, err := client.Indices.UpdateAliases(bytes.NewReader(jsonBody), client.Indices.UpdateAliases.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to update aliases: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(res.Body)

	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read alias update response body: %w", readErr)
	}

	if res.IsError() {
		return fmt.Errorf("failed to update aliases: [%s] %s", res.Status(), string(bodyBytes))
	}

	return nil
}

// deleteOldIndices finds and deletes indices matching the alias pattern, excluding the currently active one.
func deleteOldIndices(ctx context.Context, client *elasticsearch.Client, aliasName, currentIndexName string) error {
	log.Printf("Starting cleanup of old indices for alias: %s (excluding %s)", aliasName, currentIndexName)

	// Pattern to match indices for this alias
	indexPattern := fmt.Sprintf("%s_*", aliasName)

	// Use Cat Indices API to get a list of indices matching the pattern
	res, err := client.Cat.Indices(client.Cat.Indices.WithIndex(indexPattern), client.Cat.Indices.WithContext(ctx), client.Cat.Indices.WithH("index"))
	if err != nil {
		return fmt.Errorf("failed to list indices with pattern %s: %w", indexPattern, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to list indices with pattern %s: [%s] %s", indexPattern, res.Status(), string(bodyBytes))
	}

	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read indices list response body: %w", readErr)
	}

	var indicesToDelete []string
	lines := strings.Split(strings.TrimSpace(string(bodyBytes)), "\n")
	for _, line := range lines {
		indexName := strings.TrimSpace(line)
		if indexName != "" && indexName != currentIndexName {
			indicesToDelete = append(indicesToDelete, indexName)
		}
	}

	if len(indicesToDelete) == 0 {
		log.Println("No old indices found to delete.")
		return nil
	}

	log.Printf("Found %d old indices to delete: %v", len(indicesToDelete), indicesToDelete)

	// Delete the old indices
	delRes, err := client.Indices.Delete(indicesToDelete, client.Indices.Delete.WithContext(ctx), client.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil {
		return fmt.Errorf("failed to submit delete request for old indices: %w", err)
	}
	defer delRes.Body.Close()

	if delRes.IsError() {
		delBodyBytes, _ := io.ReadAll(delRes.Body)
		// Log error but don't fail the whole indexing process just because cleanup failed
		log.Printf("Warning: Failed to delete old indices: [%s] %s", delRes.Status(), string(delBodyBytes))
		return nil // Don't return error, just log it
	}

	log.Printf("Successfully deleted %d old indices.", len(indicesToDelete))
	return nil
}

// IndexMoments fetches all Moment records from the database and indexes them into a new
// Elasticsearch index, then atomically switches the alias to point to the new index.
func IndexMoments(ctx context.Context, client *elasticsearch.Client, dbService database.Service, aliasName string) error {
	if client == nil {
		return fmt.Errorf("elasticsearch client is not initialized")
	}
	if dbService == nil {
		return fmt.Errorf("database service is not initialized")
	}
	if aliasName == "" {
		return fmt.Errorf("index alias name is required")
	}

	log.Printf("Starting re-indexing process for moments with alias: %s", aliasName)

	// 1. Create a new index with a timestamp
	newIndexName := fmt.Sprintf("%s_%s", aliasName, time.Now().Format("20060102150405"))

	mapping := map[string]any{
		"settings": map[string]any{
			"number_of_shards":   3,
			"number_of_replicas": 0,
			"max_ngram_diff":     2,
			"analysis": map[string]any{
				"tokenizer": map[string]any{
					"cjk_ngram": map[string]any{
						"type":     "ngram",
						"min_gram": 2,
						"max_gram": 4,
					},
				},
				"analyzer": map[string]any{
					"ik_cjk": map[string]any{
						"tokenizer": "ik_max_word",
					},
					"cjk_ngram_analyzer": map[string]any{
						"tokenizer": "cjk_ngram",
						"filter":    []string{"lowercase"},
					},
				},
			},
		},
		"mappings": map[string]any{
			"properties": MergeTextFields([]string{"content"}),
		},
	}

	// 向量字段统一追加
	vecFields := map[string]string{
		"content_vector": "content_vector",
	}
	props := mapping["mappings"].(map[string]any)["properties"].(map[string]any)
	for name := range vecFields {
		props[name] = map[string]any{
			"type":       "dense_vector",
			"dims":       vectorDims,
			"index":      true,
			"similarity": "cosine",
		}
	}

	body, _ := json.Marshal(mapping)

	log.Printf("Creating new index: %s", newIndexName)
	res, err := client.Indices.Create(
		newIndexName,
		client.Indices.Create.WithContext(ctx),
		client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return fmt.Errorf("cannot create index %s: %w", newIndexName, err)
	}
	if res.IsError() {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Printf("Error closing response body: %v", err)
			}
		}(res.Body)
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("cannot create index %s: [%s] %s", newIndexName, res.Status(), string(bodyBytes))
	}
	log.Printf("Index %s created successfully", newIndexName)
	err = res.Body.Close()
	if err != nil {
		return err
	} // Close successful response body

	// 2. Fetch data from the database with eager loading of images and their descriptions
	log.Printf("Fetching moments from the database")
	moments, err := dbService.GetAllMoments(ctx)
	if err != nil {
		// 如果获取失败，删除刚创建的索引
		_, delErr := client.Indices.Delete([]string{newIndexName}, client.Indices.Delete.WithContext(ctx))
		if delErr != nil {
			log.Printf("Failed to delete temporary index %s after DB error: %v", newIndexName, delErr)
		}
		return fmt.Errorf("failed to fetch moments from database: %w", err)
	}
	log.Printf("Fetched %d moments from the database", len(moments))

	// 3. Index the data
	log.Printf("Indexing moments into the new index")
	if err := bulkIndexMoments(ctx, client, newIndexName, moments); err != nil {
		// 如果批量索引失败，删除刚创建的索引
		_, delErr := client.Indices.Delete([]string{newIndexName}, client.Indices.Delete.WithContext(ctx))
		if delErr != nil {
			log.Printf("Failed to delete temporary index %s after bulk index error: %v", newIndexName, delErr)
		}
		return fmt.Errorf("failed to index moments: %w", err)
	}

	// 4. Update alias to point to new index
	log.Printf("Updating alias %s to point to new index %s", aliasName, newIndexName)
	// 使用封装的 updateAlias 函数替代自定义实现
	if err := updateAlias(ctx, client, aliasName, newIndexName); err != nil {
		return fmt.Errorf("failed to update alias %s: %w", aliasName, err)
	}
	log.Printf("Alias %s updated successfully", aliasName)

	// 5. Delete old indices (run in background, log errors)
	go func() {
		if err := deleteOldIndices(context.Background(), client, aliasName, newIndexName); err != nil {
			log.Printf("Background deletion of old moment indices encountered an issue: %v", err)
		}
	}()

	log.Printf("Re-indexing process for moments with alias %s completed successfully", aliasName)
	return nil
}

// bulkIndexMoments performs bulk indexing of Moment documents.
func bulkIndexMoments(ctx context.Context, client *elasticsearch.Client, indexName string, moments []*ent.Moment) error {
	var ( // Bulk buffer and counters
		buf    bytes.Buffer
		numOps int
	)

	for _, moment := range moments {
		// Prepare meta line (action and metadata)
		meta := fmt.Sprintf(bulkIndexAction, indexName, moment.ID)
		buf.WriteString(meta)
		buf.WriteByte('\n')

		// Collect image data if available
		var images []map[string]interface{}
		if moment.Edges.MomentImages != nil {
			for _, mi := range moment.Edges.MomentImages {
				if mi.Edges.Image != nil {
					img := map[string]interface{}{
						"id":          mi.Edges.Image.ID,
						"description": mi.Edges.Image.Description,
					}
					images = append(images, img)
				}
			}
		}

		content := moment.Content
		contentVector := moment.ContentVector

		// Prepare data line (document source)
		doc := map[string]interface{}{
			"id":      moment.ID,
			"content": content,
			"images":  images,
			// Add other fields if needed for search or display
			"content_vector": contentVector,
		}
		data, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling moment %s: %v. Skipping.", moment.ID, err)
			buf.Reset() // Clear the buffer for this item
			continue
		}
		buf.Write(data)
		buf.WriteByte('\n')

		numOps++

		// Flush buffer if thresholds reached
		if buf.Len() > bulkFlushBytes || numOps >= bulkFlushItems {
			if err := flushBulkBuffer(ctx, client, &buf); err != nil {
				return err // Propagate error up
			}
			numOps = 0 // Reset counter
		}
	}

	// Flush any remaining items in the buffer
	if buf.Len() > 0 {
		if err := flushBulkBuffer(ctx, client, &buf); err != nil {
			return err
		}
	}

	// Refresh the index to make changes searchable immediately
	_, err := client.Indices.Refresh(client.Indices.Refresh.WithIndex(indexName))
	if err != nil {
		log.Printf("Warning: Failed to refresh index %s after bulk indexing: %v", indexName, err)
		// Don't fail the whole process, but log the warning
	}

	return nil
}
