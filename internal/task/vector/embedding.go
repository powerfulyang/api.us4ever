package vector

import (
	"api.us4ever/internal/ent/keep"
	"api.us4ever/internal/ent/moment"
	"api.us4ever/internal/es"
	"api.us4ever/internal/server"
	"context"
	"encoding/json"
	"log"
	"time"
)

func EmbeddingMoments(fiberServer *server.FiberServer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	db := fiberServer.DbClient
	records, err := db.Client().Moment.Query().
		Where(
			moment.ContentNEQ(""),
			moment.ContentVectorIsNil(),
		).
		All(ctx)
	if err != nil {
		return 0, err
	}

	if len(records) > 0 {
		log.Printf("EmbeddingMoments: Found %d moments to process", len(records))
	}

	for _, record := range records {
		vector, err := es.Embed(ctx, record.Content)
		if err != nil {
			// 处理错误
			log.Printf("Error embedding content for record %s: %v", record.ID, err)
			continue
		}
		// 将向量转换为 JSON 格式
		contentVector, err := json.Marshal(vector)
		if err != nil {
			// 处理错误
			log.Printf("Error marshalling vector for record %s: %v", record.ID, err)
			continue
		}
		_, err = record.Update().SetContentVector(contentVector).Save(ctx)
		if err != nil {
			// 处理错误
			log.Printf("Error updating content vector for record %s: %v", record.ID, err)
			continue
		}
	}

	// 重建 es 索引
	err = es.IndexMoments(ctx, fiberServer.EsClient, fiberServer.DbClient, fiberServer.KeepEsIndexAlias)
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func EmbeddingKeeps(fiberServer *server.FiberServer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	db := fiberServer.DbClient
	records, err := db.Client().Keep.Query().
		Where(
			keep.Or(
				keep.TitleVectorIsNil(),
				keep.SummaryVectorIsNil(),
				keep.ContentVectorIsNil(),
			),
		).
		All(ctx)
	if err != nil {
		return 0, err
	}

	if len(records) > 0 {
		log.Printf("EmbeddingKeeps: Found %d keeps to process", len(records))
	}

	for _, record := range records {
		if record.TitleVector == nil && record.Title != "" {
			vector, err := es.Embed(ctx, record.Title)
			if err != nil {
				log.Printf("Error embedding title for record %s: %v", record.ID, err)
				continue
			}
			titleVector, err := json.Marshal(vector)
			if err != nil {
				log.Printf("Error marshalling title vector for record %s: %v", record.ID, err)
				continue
			}
			_, err = record.Update().SetTitleVector(titleVector).Save(ctx)
			if err != nil {
				log.Printf("Error updating title vector for record %s: %v", record.ID, err)
				continue
			}
		}
		if record.SummaryVector == nil && record.Summary != "" {
			vector, err := es.Embed(ctx, record.Summary)
			if err != nil {
				log.Printf("Error embedding summary for record %s: %v", record.ID, err)
				continue
			}
			summaryVector, err := json.Marshal(vector)
			if err != nil {
				log.Printf("Error marshalling summary vector for record %s: %v", record.ID, err)
				continue
			}
			_, err = record.Update().SetSummaryVector(summaryVector).Save(ctx)
			if err != nil {
				log.Printf("Error updating summary vector for record %s: %v", record.ID, err)
				continue
			}
		}
		if record.ContentVector == nil && record.Content != "" {
			vector, err := es.Embed(ctx, record.Content)
			if err != nil {
				log.Printf("Error embedding content for record %s: %v", record.ID, err)
				continue
			}
			contentVector, err := json.Marshal(vector)
			if err != nil {
				log.Printf("Error marshalling content vector for record %s: %v", record.ID, err)
				continue
			}
			_, err = record.Update().SetContentVector(contentVector).Save(ctx)
			if err != nil {
				log.Printf("Error updating content vector for record %s: %v", record.ID, err)
				continue
			}
		}
	}

	// 重建 es 索引
	err = es.IndexKeeps(ctx, fiberServer.EsClient, fiberServer.DbClient, fiberServer.KeepEsIndexAlias)

	return len(records), nil
}
