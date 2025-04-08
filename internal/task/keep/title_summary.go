package keep

import (
	"context"
	"log"
	"time"

	"api.us4ever/internal/database"
	"api.us4ever/internal/dify"
	"api.us4ever/internal/ent/keep"
)

// GenerateTitleAndSummary 生成 Keep 表中缺少 title 和 summary 的记录
func GenerateTitleAndSummary(db database.Service) {
	log.Println("开始执行 GenerateKeepTitleAndSummary 任务...")

	// 创建上下文
	ctx := context.Background()

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
		log.Printf("查询 Keep 记录失败: %v", err)
		return
	}

	if len(keeps) == 0 {
		log.Println("没有需要处理的记录")
		return
	}

	log.Printf("找到 %d 条需要处理的记录", len(keeps))

	// 处理每条记录
	for _, k := range keeps {
		// 生成 title
		if k.Title == "" {
			title, err := generateTitle(k.Content)
			if err != nil {
				log.Printf("生成标题失败 (ID: %s): %v", k.ID, err)
				continue
			}

			// 更新 title
			_, err = db.Client().Keep.UpdateOne(k).
				SetTitle(title).
				SetUpdatedAt(time.Now()).
				Save(ctx)

			if err != nil {
				log.Printf("更新标题失败 (ID: %s): %v", k.ID, err)
				continue
			}

			log.Printf("成功更新标题 (ID: %s): %s", k.ID, title)
		}

		// 生成 summary
		if k.Summary == "" {
			summary, err := generateSummary(k.Content)
			if err != nil {
				log.Printf("生成摘要失败 (ID: %s): %v", k.ID, err)
				continue
			}

			// 更新 summary
			_, err = db.Client().Keep.UpdateOne(k).
				SetSummary(summary).
				SetUpdatedAt(time.Now()).
				Save(ctx)

			if err != nil {
				log.Printf("更新摘要失败 (ID: %s): %v", k.ID, err)
				continue
			}

			log.Printf("成功更新摘要 (ID: %s): %s", k.ID, summary)
		}
	}

	log.Println("GenerateKeepTitleAndSummary 任务执行完成")
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
