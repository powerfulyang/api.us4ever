package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"api.us4ever/internal/logger"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	_ "github.com/joho/godotenv/autoload"
)

var (
	nacosClient config_client.IConfigClient
	once        sync.Once
)

// NacosConfig 保存Nacos配置中心的配置
type NacosConfig struct {
	ServerAddr  string
	ServerPort  uint64
	NamespaceID string
	Group       string
	DataID      string
	Username    string
	Password    string
}

// InitNacosClient 初始化Nacos配置中心客户端
func InitNacosClient() config_client.IConfigClient {
	once.Do(func() {
		serverPort, _ := strconv.ParseUint(os.Getenv("NACOS_SERVER_PORT"), 10, 64)

		// 创建ServerConfig
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr: os.Getenv("NACOS_SERVER_ADDR"),
				Port:   serverPort,
			},
		}

		// 创建ClientConfig
		clientConfig := constant.ClientConfig{
			NamespaceId:         os.Getenv("NACOS_NAMESPACE_ID"),
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogLevel:            "error",
			Username:            os.Getenv("NACOS_USERNAME"),
			Password:            os.Getenv("NACOS_PASSWORD"),
		}

		// 创建配置客户端
		client, err := clients.CreateConfigClient(map[string]interface{}{
			"serverConfigs": serverConfigs,
			"clientConfig":  clientConfig,
		})

		if err != nil {
			configLogger.Fatal("failed to initialize Nacos config client", logger.Fields{
				"error": err.Error(),
			})
		}

		nacosClient = client
	})

	return nacosClient
}

// GetConfig 从Nacos获取配置内容
func GetConfig(dataID, group string) (string, error) {
	client := InitNacosClient()

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})

	if err != nil {
		fmt.Println("获取Nacos配置失败: ", err)
		return "", err
	}

	return content, nil
}

// ListenConfig 监听配置变化
func ListenConfig(dataID, group string, callback func(string)) error {
	client := InitNacosClient()

	err := client.ListenConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			configLogger.Info("configuration changed", logger.Fields{
				"data_id": dataId,
				"group":   group,
			})
			callback(data)
		},
	})

	if err != nil {
		return fmt.Errorf("监听Nacos配置失败: %v", err)
	}

	return nil
}

// PublishConfig 发布配置
func PublishConfig(dataID, group, content string) (bool, error) {
	client := InitNacosClient()

	success, err := client.PublishConfig(vo.ConfigParam{
		DataId:  dataID,
		Group:   group,
		Content: content,
	})

	if err != nil {
		return false, fmt.Errorf("发布Nacos配置失败: %v", err)
	}

	return success, nil
}

// LoadNacosConfig 从环境变量中加载Nacos配置
func LoadNacosConfig() NacosConfig {
	serverPort, _ := strconv.ParseUint(os.Getenv("NACOS_SERVER_PORT"), 10, 64)

	return NacosConfig{
		ServerAddr:  os.Getenv("NACOS_SERVER_ADDR"),
		ServerPort:  serverPort,
		NamespaceID: os.Getenv("NACOS_NAMESPACE_ID"),
		Group:       os.Getenv("NACOS_GROUP"),
		DataID:      os.Getenv("NACOS_DATA_ID"),
		Username:    os.Getenv("NACOS_USERNAME"),
		Password:    os.Getenv("NACOS_PASSWORD"),
	}
}
