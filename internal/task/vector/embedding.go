package vector

import (
	"context"
	"encoding/json"
	"time"

	"api.us4ever/internal/ent/keep"
	"api.us4ever/internal/ent/moment"
	"api.us4ever/internal/es"
	"api.us4ever/internal/logger"
	"api.us4ever/internal/server"
	"go.uber.org/zap"
)

var (
	embeddingLogger *logger.Logger
)

func init() {
	var err error
	embeddingLogger, err = logger.New("embedding")
	if err != nil {
		panic("failed to initialize embedding logger: " + err.Error())
	}
}

func EmbeddingMoments(fiberServer *server.FiberServer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
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
		embeddingLogger.Info("found moments to process for embedding",
			zap.Int("count", len(records)),
		)
	}

	handledCount := 0

	for _, record := range records {
		vector, err := es.Embed(ctx, record.Content)
		if err != nil {
			// Handle error
			embeddingLogger.Error("error embedding content for moment record",
				zap.String("record_id", record.ID),
				zap.Error(err),
			)
			continue
		}
		// Convert vector to JSON format
		contentVector, err := json.Marshal(vector)
		if err != nil {
			// Handle error
			embeddingLogger.Error("error marshalling vector for moment record",
				zap.String("record_id", record.ID),
				zap.Error(err),
			)
			continue
		}
		_, err = record.Update().SetContentVector(contentVector).Save(ctx)
		if err != nil {
			// Handle error
			embeddingLogger.Error("error updating content vector for moment record",
				zap.String("record_id", record.ID),
				zap.Error(err),
			)
			continue
		}
		handledCount++
	}

	// 重建 es 索引
	if handledCount > 0 {
		err = es.IndexMoments(ctx, fiberServer.EsClient, fiberServer.DbClient, fiberServer.MomentEsIndexAlias)
	}
	if err != nil {
		return 0, err
	}
	return handledCount, nil
}

func EmbeddingKeeps(fiberServer *server.FiberServer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
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
		embeddingLogger.Info("found keeps to process for embedding",
			zap.Int("count", len(records)),
		)
	}

	handledCount := 0

	for _, record := range records {
		if record.TitleVector == nil && record.Title != "" {
			vector, err := es.Embed(ctx, record.Title)
			if err != nil {
				embeddingLogger.Error("error embedding title for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
			titleVector, err := json.Marshal(vector)
			if err != nil {
				embeddingLogger.Error("error marshalling title vector for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
			_, err = record.Update().SetTitleVector(titleVector).Save(ctx)
			if err != nil {
				embeddingLogger.Error("error updating title vector for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
		}
		if record.SummaryVector == nil && record.Summary != "" {
			vector, err := es.Embed(ctx, record.Summary)
			if err != nil {
				embeddingLogger.Error("error embedding summary for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
			summaryVector, err := json.Marshal(vector)
			if err != nil {
				embeddingLogger.Error("error marshalling summary vector for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
			_, err = record.Update().SetSummaryVector(summaryVector).Save(ctx)
			if err != nil {
				embeddingLogger.Error("error updating summary vector for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
		}
		if record.ContentVector == nil && record.Content != "" {
			vector, err := es.Embed(ctx, record.Content)
			if err != nil {
				embeddingLogger.Error("error embedding content for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
			contentVector, err := json.Marshal(vector)
			if err != nil {
				embeddingLogger.Error("error marshalling content vector for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
			_, err = record.Update().SetContentVector(contentVector).Save(ctx)
			if err != nil {
				embeddingLogger.Error("error updating content vector for keep record",
					zap.String("record_id", record.ID),
					zap.Error(err),
				)
				continue
			}
		}
		handledCount++
	}

	// 重建 es 索引
	if handledCount > 0 {
		err = es.IndexKeeps(ctx, fiberServer.EsClient, fiberServer.DbClient, fiberServer.KeepEsIndexAlias)
	}

	if err != nil {
		return 0, err
	}

	return handledCount, nil
}
