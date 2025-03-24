package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// AppConfig 应用配置结构体
type AppConfig struct {
	AppName  string       `json:"app_name"`
	AppEnv   string       `json:"app_env"`
	Server   ServerConfig `json:"server"`
	Database DBConfig     `json:"database"`
	Redis    RedisConfig  `json:"redis,omitempty"`
	// 添加其他配置项...
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int `json:"port"`
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password,omitempty"`
	DB       int    `json:"db"`
}

var (
	appConfig *AppConfig
	configMux sync.RWMutex
)

// LoadConfig 加载配置
func LoadConfig() (*AppConfig, error) {
	configMux.Lock()
	defer configMux.Unlock()

	if appConfig != nil {
		return appConfig, nil
	}

	// 从Nacos加载配置
	nacosConfig := LoadNacosConfig()
	configContent, err := GetConfig(nacosConfig.DataID, nacosConfig.Group)

	// 如果无法从Nacos加载，则使用环境变量
	if err != nil {
		log.Printf("从Nacos加载配置失败，使用环境变量。")
		return loadConfigFromEnv(), nil
	}

	// 解析Nacos配置
	config := &AppConfig{}
	if err := json.Unmarshal([]byte(configContent), config); err != nil {
		log.Printf("解析Nacos配置失败，使用环境变量: %v", err)
		return loadConfigFromEnv(), nil
	}

	appConfig = config

	// 设置监听配置变化
	setupConfigListener(nacosConfig.DataID, nacosConfig.Group)

	return appConfig, nil
}

// loadConfigFromEnv 从环境变量加载配置
func loadConfigFromEnv() *AppConfig {
	config := &AppConfig{
		AppName: "api.us4ever",
		AppEnv:  os.Getenv("APP_ENV"),
		Server: ServerConfig{
			Port: 8080, // 默认值，应该从环境变量获取
		},
		Database: DBConfig{
			Host:     os.Getenv("BLUEPRINT_DB_HOST"),
			Port:     5432, // 应该从环境变量获取
			Database: os.Getenv("BLUEPRINT_DB_DATABASE"),
			Username: os.Getenv("BLUEPRINT_DB_USERNAME"),
			Password: os.Getenv("BLUEPRINT_DB_PASSWORD"),
			Schema:   os.Getenv("BLUEPRINT_DB_SCHEMA"),
		},
	}

	return config
}

// GetConfig 获取配置
func GetAppConfig() *AppConfig {
	configMux.RLock()
	defer configMux.RUnlock()

	if appConfig == nil {
		configMux.RUnlock()
		_, err := LoadConfig()
		if err != nil {
			log.Printf("加载配置失败: %v", err)
		}
		configMux.RLock()
	}

	return appConfig
}

// setupConfigListener 设置配置变更监听
func setupConfigListener(dataID, group string) {
	err := ListenConfig(dataID, group, func(content string) {
		configMux.Lock()
		defer configMux.Unlock()

		newConfig := &AppConfig{}
		if err := json.Unmarshal([]byte(content), newConfig); err != nil {
			log.Printf("解析更新的配置失败: %v", err)
			return
		}

		// 更新配置
		appConfig = newConfig
		log.Println("配置已更新")
	})

	if err != nil {
		log.Printf("设置配置监听失败: %v", err)
	}
}
