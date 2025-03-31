package config

import (
	"fmt"
)

// GetDSN 获取数据库连接字符串
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s",
		c.Username, c.Password, c.Host, c.Port, c.Database, c.Schema)
}

func LoadDatabaseConfig() (*DBConfig, error) {
	// 使用现有的配置包
	appConfig, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load app config: %v", err)
	}

	// 转换配置结构
	dbConfig := appConfig.Database

	// 验证数据库配置
	if err := dbConfig.Validate(); err != nil {
		return nil, err
	}

	return &dbConfig, nil
}
