//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"api.us4ever/internal/config"
	"api.us4ever/internal/es"

	"api.us4ever/internal/logger"
)

var (
	searchTestLogger *logger.Logger
)

func init() {
	var err error
	searchTestLogger, err = logger.New("search-test")
	if err != nil {
		panic("failed to initialize search-test logger: " + err.Error())
	}
}

func TestSearchKeeps(t *testing.T) {
	ctx := context.Background()

	appConfig := config.GetAppConfig()
	if appConfig == nil {
		searchTestLogger.Fatal("failed to load application config")
	}

	client, err := es.NewClient(appConfig.ES)
	if err != nil {
		searchTestLogger.Fatal("failed to initialize Elasticsearch client", logger.Fields{
			"error": err.Error(),
		})
	}

	indexAlias := fmt.Sprintf("%s-keeps", strings.ToLower(strings.ReplaceAll(appConfig.AppName, " ", "-")))

	searchQuery := "测试"

	searchTestLogger.Info("performing search test", logger.Fields{
		"index_alias":  indexAlias,
		"search_query": searchQuery,
	})
	results, err := es.SearchKeeps(ctx, client, indexAlias, searchQuery)

	total := results.Hits.Total.Value

	// 7. 输出结果
	fmt.Printf("搜索结果总数: %d\n", total)
}
