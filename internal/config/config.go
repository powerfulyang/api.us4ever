package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"api.us4ever/internal/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	configLogger    *logger.Logger
)

func init() {
	var err error
	configLogger, err = logger.New("config")
	if err != nil {
		panic("failed to initialize config logger: " + err.Error())
	}
}

// RegisterChangeCallback registers a callback function to be called when config changes
func RegisterChangeCallback(callback ChangeCallback) {
	callbacksMutex.Lock()
	defer callbacksMutex.Unlock()
	changeCallbacks = append(changeCallbacks, callback)
}

// LoadConfig loads the application configuration
func LoadConfig() (*AppConfig, error) {
	configMux.Lock()
	defer configMux.Unlock()

	// Return cached config if available
	if appConfig != nil {
		return appConfig, nil
	}

	// Load environment variables if not in container environment
	if err := loadEnvironmentFile(); err != nil {
		configLogger.Warn("failed to load .env file", zap.Error(err))
	}

	// Load configuration from Nacos
	config, err := loadFromNacos()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration from Nacos: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Cache the configuration
	appConfig = config

	// Setup configuration change listener
	nacosConfig := LoadNacosConfig()
	setupConfigListener(nacosConfig.DataID, nacosConfig.Group)

	configLogger.Info("configuration loaded successfully")
	return appConfig, nil
}

// loadEnvironmentFile loads the .env file if it exists
func loadEnvironmentFile() error {
	nacosServerAddr := os.Getenv("NACOS_SERVER_ADDR")
	if nacosServerAddr != "" {
		// Skip loading .env file in container environment
		return nil
	}

	// Try to load .env file from current directory first
	if err := godotenv.Load(".env"); err == nil {
		return nil
	}

	if err := godotenv.Load("../.env"); err == nil {
		return nil
	}

	// Try to load .env file from current directory first
	if err := godotenv.Load("../../.env"); err == nil {
		return nil
	}

	if err := godotenv.Load("../../../.env"); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	return nil
}

// loadFromNacos loads configuration from Nacos
func loadFromNacos() (*AppConfig, error) {
	nacosConfig := LoadNacosConfig()
	configContent, err := GetConfig(nacosConfig.DataID, nacosConfig.Group)
	if err != nil {
		return nil, fmt.Errorf("failed to get config from Nacos: %w", err)
	}

	config := &AppConfig{}
	if err := json.Unmarshal([]byte(configContent), config); err != nil {
		return nil, fmt.Errorf("failed to parse Nacos configuration: %w", err)
	}

	return config, nil
}

// GetAppConfig returns the application configuration
// This function should only be used after ensuring LoadConfig() has been called successfully
func GetAppConfig() *AppConfig {
	config, err := LoadConfig()
	if err != nil {
		configLogger.Error("failed to load configuration", zap.Error(err))
		return nil
	}
	return config
}

// MustGetAppConfig returns the application configuration or panics if it fails
// Use this only during application startup when configuration is critical
func MustGetAppConfig() *AppConfig {
	config, err := LoadConfig()
	if err != nil {
		configLogger.Fatal("failed to load configuration", zap.Error(err))
	}
	return config
}

// setupConfigListener sets up configuration change listener
func setupConfigListener(dataID, group string) {
	err := ListenConfig(dataID, group, func(content string) {
		configMux.Lock()
		defer configMux.Unlock()

		newConfig := &AppConfig{}
		if err := json.Unmarshal([]byte(content), newConfig); err != nil {
			configLogger.Error("failed to parse updated configuration", zap.Error(err))
			return
		}

		// Validate the new configuration
		if err := newConfig.Validate(); err != nil {
			configLogger.Error("updated configuration validation failed", zap.Error(err))
			return
		}

		// Update the cached configuration
		appConfig = newConfig
		configLogger.Info("configuration updated successfully")

		// Notify all registered callbacks
		notifyConfigChange(newConfig)
	})

	if err != nil {
		configLogger.Error("failed to setup configuration listener", zap.Error(err))
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

// Validate validates the database configuration
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
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535, got %d", c.Port)
	}
	if c.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Schema == "" {
		return fmt.Errorf("database schema is required")
	}
	return nil
}

// Validate validates the server configuration
func (c *ServerConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535, got %d", c.Port)
	}
	return nil
}

// Validate validates the application configuration
func (c *AppConfig) Validate() error {
	if c.AppName == "" {
		return fmt.Errorf("app name is required")
	}
	if c.AppEnv == "" {
		return fmt.Errorf("app environment is required")
	}

	if err := c.Server.Validate(); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := c.Database.Validate(); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}

	return nil
}
