package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
	"go.uber.org/zap"
)

var (
	telegramSyncLogger *logger.Logger
)

func init() {
	var err error
	telegramSyncLogger, err = logger.New("telegram-sync")
	if err != nil {
		panic("failed to initialize telegram-sync logger: " + err.Error())
	}
}

func TriggerSyncTelegram() (int, error) {
	// 向 https://us4ever.com/api/internal/sync/telegram/emt_channel 发送 GET 请求
	appConfig := config.GetAppConfig()
	url := appConfig.Telegram.SyncURL
	if url == "" {
		return 0, fmt.Errorf("telegram Sync URL is not configured")
	}

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			telegramSyncLogger.Error("failed to close response body",
				zap.Error(err),
			)
		}
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	var result struct {
		Success bool `json:"success"`
		Count   int  `json:"count"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("解析响应失败: %v", err)
	}

	return result.Count, nil
}
