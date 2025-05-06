package telegram

import (
	"api.us4ever/internal/config"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
			log.Printf("关闭响应体失败: %v", err)
		}
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return 1, nil
}
