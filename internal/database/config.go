package database

import (
	"fmt"

	"api.us4ever/internal/config"
)

// DBConfig 数据库配置结构
type DBConfig struct {
	Database string
	Password string
	Username string
	Port     int
	Host     string
	Schema   string
}

// GetDSN 获取数据库连接字符串
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Schema)
}

// LoadConfig 从应用配置加载数据库配置
func LoadConfig() (*DBConfig, error) {
	// 使用现有的配置包
	appConfig, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load app config: %v", err)
	}

	// 转换配置结构
	dbConfig := &DBConfig{
		Host:     appConfig.Database.Host,
		Port:     appConfig.Database.Port,
		Database: appConfig.Database.Database,
		Username: appConfig.Database.Username,
		Password: appConfig.Database.Password,
		Schema:   appConfig.Database.Schema,
	}

	// 验证数据库配置
	if err := dbConfig.Validate(); err != nil {
		return nil, err
	}

	return dbConfig, nil
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
