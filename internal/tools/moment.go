package tools

import (
	"api.us4ever/internal/database"
	"api.us4ever/internal/es"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
)

// ImportMomentsFromCSV 从 CSV 文件导入数据到 moment 表
func ImportMomentsFromCSV(filePath string) error {
	// 创建上下文
	ctx := context.Background()
	// 初始化数据库服务
	db, err := database.New()
	// 打开 CSV 文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开 CSV 文件失败: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("关闭 CSV 文件失败: %v\n", err)
		}
		err = db.Close()
		if err != nil {
			fmt.Printf("关闭数据库连接失败: %v\n", err)
		}
	}(file)

	// 创建 CSV reader
	reader := csv.NewReader(file)

	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("读取 CSV 表头失败: %w", err)
	}

	// 验证 content 字段
	contentIndex := -1
	for i, header := range headers {
		if header == "content" {
			contentIndex = i
			break
		}
	}
	if contentIndex == -1 {
		return fmt.Errorf("CSV 文件缺少 content 字段")
	}

	userId := db.Client().User.Query().FirstIDX(ctx)

	// 读取并处理每一行数据
	for {
		record, err := reader.Read()
		if err != nil {
			break // 文件结束
		}

		if contentIndex >= len(record) {
			continue
		}

		content := record[contentIndex]
		if content == "" {
			continue
		}
		contentVector, _ := es.Embed(ctx, content)
		vectorJSON, _ := json.Marshal(contentVector)

		// 创建新的 moment
		now := time.Now()

		_, err = db.Client().Moment.Create().
			SetID(uuid.New().String()).
			SetContent(content).
			SetContentVector(vectorJSON).
			SetCategory("default").
			SetIsPublic(true).
			SetLikes(0).
			SetViews(0).
			SetTags(json.RawMessage("[]")).
			SetExtraData(json.RawMessage("{}")).
			SetOwnerId(userId).
			SetCreatedAt(now).
			SetUpdatedAt(now).
			Save(ctx)

		if err != nil {
			return fmt.Errorf("保存 moment 失败: %w", err)
		}

		log.Printf("成功导入 moment: %s\n", content)
	}

	return nil
}
