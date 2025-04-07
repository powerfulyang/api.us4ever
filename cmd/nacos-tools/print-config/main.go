package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"api.us4ever/internal/config"
)

func main() {
	// 获取配置
	content := config.GetAppConfig()

	// 格式化并打印配置
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		// 如果序列化失败，直接打印原始内容
		fmt.Printf("%+v\n", content)
		return
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonBytes, "", "  "); err != nil {
		// 如果不是有效的JSON，直接打印原始内容
		fmt.Printf("%+v\n", content)
	} else {
		fmt.Println(prettyJSON.String())
	}
}
