//go:build e2e

package e2e

import (
	"api.us4ever/internal/config"
	"api.us4ever/internal/es"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestSearchKeeps(t *testing.T) {
	ctx := context.Background()

	appConfig := config.GetAppConfig()
	if appConfig == nil {
		log.Fatalf("无法加载应用配置")
	}

	client, err := es.NewClient(appConfig.ES)
	if err != nil {
		log.Fatalf("初始化 Elasticsearch 客户端失败: %v", err)
	}

	indexAlias := fmt.Sprintf("%s-keeps", strings.ToLower(strings.ReplaceAll(appConfig.AppName, " ", "-")))

	searchQuery := "测试"

	log.Printf("使用索引别名 '%s' 搜索关键词 '%s'", indexAlias, searchQuery)
	results, err := es.SearchKeeps(ctx, client, indexAlias, searchQuery)

	total := results.Hits.Total.Value

	// 7. 输出结果
	fmt.Printf("搜索结果总数: %d\n", total)
}
