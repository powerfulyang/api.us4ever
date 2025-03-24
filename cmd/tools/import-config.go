package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"api.us4ever/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// 定义命令行参数
	filePath := flag.String("file", "", "配置文件路径")
	dataID := flag.String("dataId", "", "Nacos配置的DataID (默认使用环境变量NACOS_DATA_ID)")
	group := flag.String("group", "", "Nacos配置的Group (默认使用环境变量NACOS_GROUP)")

	flag.Parse()

	// 验证参数
	if *filePath == "" {
		fmt.Println("错误: 必须指定配置文件路径")
		fmt.Println("使用方法: import-config -file=<配置文件路径> [-dataId=<DataID>] [-group=<Group>]")
		os.Exit(1)
	}

	// 如果未指定DataID，则使用环境变量
	if *dataID == "" {
		*dataID = os.Getenv("NACOS_DATA_ID")
		if *dataID == "" {
			fmt.Println("错误: 未指定DataID，环境变量NACOS_DATA_ID也未设置")
			os.Exit(1)
		}
	}

	// 如果未指定Group，则使用环境变量
	if *group == "" {
		*group = os.Getenv("NACOS_GROUP")
		if *group == "" {
			*group = "DEFAULT_GROUP" // 默认值
		}
	}

	// 初始化Nacos客户端
	config.InitNacosClient()

	// 导入配置
	err := config.ImportConfigToNacos(*filePath, *dataID, *group)
	if err != nil {
		log.Fatalf("导入配置失败: %v", err)
	}

	fmt.Println("配置导入成功!")
}
