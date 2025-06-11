package keep

import (
	"context"
	"fmt"
	"time"

	"api.us4ever/internal/server"

	"api.us4ever/internal/dify"
	"api.us4ever/internal/ent/keep"
	"api.us4ever/internal/logger"
)

var (
	titleSummaryLogger *logger.Logger
)

func init() {
	var err error
	titleSummaryLogger, err = logger.New("title-summary")
	if err != nil {
		panic("failed to initialize title-summary logger: " + err.Error())
	}
}

// GenerateTitleAndSummary 生成 Keep 表中缺少 title 和 summary 的记录
func GenerateTitleAndSummary(fiberServer *server.FiberServer) (int, error) {
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	db := fiberServer.DbClient
	// 查询缺少 title 或 summary 的记录
	keeps, err := db.Client().Keep.Query().
		Where(
			keep.Or(
				keep.TitleEQ(""),
				keep.SummaryEQ(""),
			),
		).
		Limit(1). // 每次处理 1 条记录
		All(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to query keeps: %v", err)
	}

	if len(keeps) > 0 {
		titleSummaryLogger.Info("found keeps to process for title/summary generation", logger.Fields{
			"count": len(keeps),
		})
	}

	// 处理每条记录
	for _, k := range keeps {
		// 生成 title
		if k.Title == "" {
			title, err := generateTitle(k.Content)
			if err != nil {
				titleSummaryLogger.Error("error generating title", logger.Fields{
					"keep_id": k.ID,
					"error":   err.Error(),
				})
				continue
			}

			// 更新 title
			_, err = k.Update().
				SetTitle(title).
				SetUpdatedAt(time.Now()).
				Save(ctx)

			if err != nil {
				titleSummaryLogger.Error("error updating title", logger.Fields{
					"keep_id": k.ID,
					"error":   err.Error(),
				})
				continue
			}
		}

		// 生成 summary
		if k.Summary == "" {
			summary, err := generateSummary(k.Content)
			if err != nil {
				titleSummaryLogger.Error("error generating summary", logger.Fields{
					"keep_id": k.ID,
					"error":   err.Error(),
				})
				continue
			}

			// 更新 summary
			_, err = k.Update().
				SetSummary(summary).
				SetUpdatedAt(time.Now()).
				Save(ctx)

			if err != nil {
				titleSummaryLogger.Error("error updating summary", logger.Fields{
					"keep_id": k.ID,
					"error":   err.Error(),
				})
				continue
			}
		}
	}

	return len(keeps), nil
}

// generateTitle 使用 Dify 生成标题
func generateTitle(content string) (string, error) {
	req := &dify.WorkflowRequest{
		Inputs: dify.WorkflowInput{
			Action:  dify.ActionTitle,
			Content: content,
		},
		ResponseMode: dify.ResponseModeBlocking,
	}

	result, err := dify.CallWorkflow(req)
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

// generateSummary 使用 Dify 生成摘要
func generateSummary(content string) (string, error) {
	req := &dify.WorkflowRequest{
		Inputs: dify.WorkflowInput{
			Action:  dify.ActionContent,
			Content: content,
		},
		ResponseMode: dify.ResponseModeBlocking,
	}

	result, err := dify.CallWorkflow(req)
	if err != nil {
		return "", err
	}

	return result.Message, nil
}
