package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	// 检查文件路径
	if *filePath == "" {
		log.Fatal("错误: 请指定配置文件路径，使用 -file 参数")
	}

	// 如果未指定DataID，则使用环境变量
	if *dataID == "" {
		*dataID = os.Getenv("NACOS_DATA_ID")
		if *dataID == "" {
			log.Fatal("错误: 未指定DataID，环境变量NACOS_DATA_ID也未设置")
		}
	}

	// 如果未指定Group，则使用环境变量
	if *group == "" {
		*group = os.Getenv("NACOS_GROUP")
		if *group == "" {
			*group = "DEFAULT_GROUP"
		}
	}

	// 读取配置文件
	content, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 初始化 Nacos 客户端
	config.InitNacosClient()

	// 发布配置到 Nacos
	success, err := config.PublishConfig(*dataID, *group, string(content))
	if err != nil {
		log.Fatalf("发布配置失败: %v", err)
	}

	if success {
		fmt.Printf("配置已成功发布到 Nacos (DataID: %s, Group: %s)\n", *dataID, *group)
	} else {
		fmt.Println("发布配置失败")
	}
}
