package main

import (
	"flag"
	"fmt"
	"os"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
)

var (
	nacosToolsLogger *logger.Logger
)

func init() {
	var err error
	nacosToolsLogger, err = logger.New("nacos-tools")
	if err != nil {
		panic("failed to initialize nacos-tools logger: " + err.Error())
	}
}

func main() {
	// 定义命令行参数
	filePath := flag.String("file", "", "配置文件路径")
	dataID := flag.String("dataId", "", "Nacos配置的DataID (默认使用环境变量NACOS_DATA_ID)")
	group := flag.String("group", "", "Nacos配置的Group (默认使用环境变量NACOS_GROUP)")
	flag.Parse()

	// 检查文件路径
	if *filePath == "" {
		nacosToolsLogger.Fatal("error: please specify config file path using -file parameter")
	}

	// 如果未指定DataID，则使用环境变量
	if *dataID == "" {
		*dataID = os.Getenv("NACOS_DATA_ID")
		if *dataID == "" {
			nacosToolsLogger.Fatal("error: DataID not specified and NACOS_DATA_ID environment variable not set")
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
	content, err := os.ReadFile(*filePath)
	if err != nil {
		nacosToolsLogger.Fatalw("failed to read config file",
			"file_path", *filePath,
			"error", err,
		)
	}

	// 初始化 Nacos 客户端
	config.InitNacosClient()

	// 发布配置到 Nacos
	success, err := config.PublishConfig(*dataID, *group, string(content))
	if err != nil {
		nacosToolsLogger.Fatalw("failed to publish config",
			"data_id", *dataID,
			"group", *group,
			"error", err,
		)
	}

	if success {
		nacosToolsLogger.Infow("config published successfully to Nacos",
			"data_id", *dataID,
			"group", *group,
		)
		fmt.Printf("配置已成功发布到 Nacos (DataID: %s, Group: %s)\n", *dataID, *group)
	} else {
		nacosToolsLogger.Errorw("failed to publish config to Nacos",
			"data_id", *dataID,
			"group", *group,
		)
		fmt.Println("发布配置失败")
	}
}
