package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ImportConfigToNacos 将本地配置文件导入到Nacos
func ImportConfigToNacos(filePath, dataID, group string) error {
	// 读取本地配置文件
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 验证JSON格式是否正确
	var jsonData interface{}
	if err := json.Unmarshal(content, &jsonData); err != nil {
		return fmt.Errorf("配置文件不是有效的JSON格式: %v", err)
	}

	// 发布配置到Nacos
	success, err := PublishConfig(dataID, group, string(content))
	if err != nil {
		return fmt.Errorf("发布配置到Nacos失败: %v", err)
	}

	if !success {
		return fmt.Errorf("发布配置到Nacos未成功")
	}

	log.Printf("成功将配置导入到Nacos: dataID=%s, group=%s", dataID, group)
	return nil
}
