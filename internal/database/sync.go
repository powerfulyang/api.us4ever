package database

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
)

// SyncSchema 从数据库同步结构到 ENT schema
func SyncSchema() error {
	// 确保目录存在
	schemaDir := "ent/schema"
	if err := os.MkdirAll(schemaDir, 0755); err != nil {
		return fmt.Errorf("failed to create schema directory: %v", err)
	}

	// 从 config 包获取配置
	dbConfig, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("无法加载应用配置: %v", err)
	}

	// 构建 DSN
	dsn := dbConfig.GetDSN()

	log.Printf("使用 DSN: %s", dsn)

	cmd := exec.Command("go", "run", "-mod=mod", "ariga.io/entimport/cmd/entimport", "-dsn", dsn)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run entimport: %v", err)
	}

	log.Printf("ENT schema generated successfully")
	return nil
}
