package telegram

import (
	"io"
	"log"
	"net/http"
)

func TriggerSyncTelegram() {
	// 向 https://us4ever.com/api/internal/sync/telegram/emt_channel发送 GET 请求
	url := "http://us4ever.com:3000/api/internal/sync/telegram/emt_channel"

	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("关闭响应体失败: %v", err)
		}
	}(resp.Body)

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("请求失败，状态码: %d", resp.StatusCode)
		return
	}

	log.Println("请求成功，已触发 Telegram 同步")
}
