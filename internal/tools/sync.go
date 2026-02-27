package tools

import (
	"fmt"
	"os"
	"os/exec"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
	_ "github.com/lib/pq"
)

var (
	syncLogger *logger.Logger
)

func init() {
	var err error
	syncLogger, err = logger.New("sync")
	if err != nil {
		panic("failed to initialize sync logger: " + err.Error())
	}
}

// SyncSchema 从数据库同步结构到 ENT schema
func SyncSchema() error {
	// 确保目录存在
	schemaDir := "internal/ent/schema"
	// 先清空目录
	if err := os.RemoveAll(schemaDir); err != nil {
		return fmt.Errorf("failed to remove existing schema directory: %v", err)
	}

	if err := os.MkdirAll(schemaDir, 0755); err != nil {
		return fmt.Errorf("failed to create schema directory: %v", err)
	}

	// 从 config 包获取配置
	dbConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		return fmt.Errorf("无法加载应用配置: %v", err)
	}

	// 构建 DSN
	dsn := dbConfig.GetDSN()

	syncLogger.Infow("using DSN for schema sync", "dsn", dsn)

	cmd := exec.Command(
		"go", "run",
		"-mod=mod", "github.com/powerfulyang/entimport/cmd/entimport",
		"-dsn", dsn,
		"-schema-path", "./internal/ent/schema",
		"--exclude-tables", "_prisma_migrations",
	)
	cmd.Dir = "."
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run entimport: %v", err)
	}

	syncLogger.Info("ENT schema generated successfully")
	return nil
}
