package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type ChangeCallback func(newConfig *AppConfig)

// AppConfig 应用配置结构体
type AppConfig struct {
	AppName   string          `json:"app_name"`
	AppEnv    string          `json:"app_env"`
	Server    ServerConfig    `json:"server"`
	Database  DBConfig        `json:"database"`
	Redis     RedisConfig     `json:"redis,omitempty"`
	Dify      DifyConfig      `json:"dify,omitempty"`
	ES        ESConfig        `json:"es,omitempty"`
	OCR       OCRConfig       `json:"ocr,omitempty"`
	Telegram  TelegramConfig  `json:"telegram,omitempty"`
	Embedding EmbeddingConfig `json:"embedding,omitempty"`
	// 添加其他配置项...
}

type EmbeddingConfig struct {
	Endpoint string `json:"endpoint"`
}

type TelegramConfig struct {
	SyncURL string `json:"sync_url"`
}

type OCRConfig struct {
	Endpoint string `json:"endpoint"`
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

// DifyConfig Dify配置
type DifyConfig struct {
	Endpoint string `json:"endpoint"`
	ApiKey   string `json:"api_key"`
}

// ESConfig Elasticsearch 配置
type ESConfig struct {
	Addresses []string `json:"addresses"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
}

var (
	appConfig       *AppConfig
	configMux       sync.RWMutex
	changeCallbacks []ChangeCallback
	callbacksMutex  sync.RWMutex
)

// RegisterChangeCallback registers a callback function to be called when config changes
func RegisterChangeCallback(callback ChangeCallback) {
	callbacksMutex.Lock()
	defer callbacksMutex.Unlock()
	changeCallbacks = append(changeCallbacks, callback)
}

// LoadConfig 加载配置
func LoadConfig() (*AppConfig, error) {
	configMux.Lock()
	defer configMux.Unlock()

	if appConfig != nil {
		return appConfig, nil
	}

	nacosServerAddr := os.Getenv("NACOS_SERVER_ADDR")

	if nacosServerAddr == "" {
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("⚠️ 加载 .env 文件失败，程序即将退出")
		}
	}

	// 从Nacos加载配置
	nacosConfig := LoadNacosConfig()
	configContent, err := GetConfig(nacosConfig.DataID, nacosConfig.Group)

	// 如果无法从Nacos加载，则报错
	if err != nil {
		return nil, fmt.Errorf("无法加载Nacos配置: %v", err)
	}

	// 解析Nacos配置
	config := &AppConfig{}
	if err := json.Unmarshal([]byte(configContent), config); err != nil {
		return nil, fmt.Errorf("解析Nacos配置失败: %v", err)
	}

	appConfig = config

	// 设置监听配置变化
	setupConfigListener(nacosConfig.DataID, nacosConfig.Group)

	return appConfig, nil
}

// GetAppConfig GetConfig 获取配置
func GetAppConfig() *AppConfig {
	config, _ := LoadConfig()
	return config
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

		// 调用所有注册的回调函数
		notifyConfigChange(newConfig)
	})

	if err != nil {
		log.Printf("设置配置监听失败: %v", err)
	}
}

// notifyConfigChange calls all registered callback functions with the new config
func notifyConfigChange(newConfig *AppConfig) {
	callbacksMutex.RLock()
	defer callbacksMutex.RUnlock()

	for _, callback := range changeCallbacks {
		go callback(newConfig)
	}
}

// Validate 验证数据库配置
func (c *DBConfig) Validate() error {
	if c.Database == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Password == "" {
		return fmt.Errorf("database password is required")
	}
	if c.Username == "" {
		return fmt.Errorf("database username is required")
	}
	if c.Port == 0 {
		return fmt.Errorf("database port is required")
	}
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Schema == "" {
		return fmt.Errorf("database schema is required")
	}
	return nil
}
