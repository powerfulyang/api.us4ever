package main

import (
	"fmt"
	"log"
	"os"

	"api.us4ever/internal/tools"
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
	default:
		log.Printf("未知命令: %s", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  go run ./cmd/db-tools sync               # 从数据库同步结构")
}

func syncSchema() {
	log.Println("正在从数据库同步结构...")

	if err := tools.SyncSchema(); err != nil {
		log.Fatalf("同步数据库结构失败: %v", err)
	}

	log.Println("数据库结构同步成功！")
}
