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
	case "import-moments":
		if len(os.Args) < 3 {
			log.Fatal("请指定 CSV 文件路径")
		}
		importMoments(os.Args[2])
	default:
		log.Printf("未知命令: %s", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  go run ./cmd/db-tools sync               # 从数据库同步结构")
	fmt.Println("  go run ./cmd/db-tools import-moments <csv文件路径>  # 从 CSV 导入数据到 moment 表")
}

func syncSchema() {
	log.Println("正在从数据库同步结构...")

	if err := tools.SyncSchema(); err != nil {
		log.Fatalf("同步数据库结构失败: %v", err)
	}

	log.Println("数据库结构同步成功！")
}

func importMoments(csvPath string) {
	log.Println("正在从 CSV 导入数据到 moment 表...")

	if err := tools.ImportMomentsFromCSV(csvPath); err != nil {
		log.Fatalf("导入数据失败: %v", err)
	}

	log.Println("数据导入成功！")
}
