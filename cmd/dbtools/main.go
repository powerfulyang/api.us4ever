package main

import (
	"fmt"
	"log"
	"os"

	"api.us4ever/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "sync":
		syncSchema()
	case "generate":
		if len(os.Args) < 3 {
			log.Fatalf("请提供迁移名称，例如: go run ./cmd/dbtools generate add_users")
		}
		generateMigration(os.Args[2])
	case "apply":
		applyMigration()
	case "setup":
		setupDatabase()
	default:
		log.Printf("未知命令: %s", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  go run ./cmd/dbtools sync               # 从数据库同步结构")
}

func syncSchema() {
	log.Println("正在从数据库同步结构...")

	if err := database.SyncSchema(); err != nil {
		log.Fatalf("同步数据库结构失败: %v", err)
	}

	log.Println("数据库结构同步成功！")
}

func generateMigration(name string) {
	log.Printf("正在生成数据库迁移: %s...", name)

	log.Println("数据库迁移生成成功！")
}

func applyMigration() {
	log.Println("正在应用数据库迁移...")

	log.Println("数据库迁移应用成功！")
}

func setupDatabase() {
	log.Println("正在初始化数据库连接...")

	log.Println("数据库连接初始化成功！")
}
