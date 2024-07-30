package initiliza

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"user/user_buf/global"
)

// 初始化配置文件
func InitConfig() {
	v := viper.New()
	v.SetConfigFile("./config/dev_config.yaml")

	if err := v.ReadInConfig(); err != nil {
		zap.S().Panic("配置文件读取失败")
	}

	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		zap.S().Panic("配置文件获取失败")
	}

	InitNacos()

	log.Println("配置文件", global.ServerConfig)

}

// 初始化nacos
func InitNacos() {
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.ServerConfig.Nacos.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "./tmp/nacos/log",
		CacheDir:            "./tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.ServerConfig.Nacos.ServerIp,
			ContextPath: "/nacos",
			Port:        global.ServerConfig.Nacos.ServerPort,
			Scheme:      "http",
		},
	}
	client, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	content, _ := client.GetConfig(vo.ConfigParam{
		DataId: global.ServerConfig.Nacos.DataId,
		Group:  global.ServerConfig.Nacos.Group})

	json.Unmarshal([]byte(content), &global.ServerConfig)

}
