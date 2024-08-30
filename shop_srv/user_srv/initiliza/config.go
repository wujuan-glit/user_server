package initiliza

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"user/user_srv/global"
)

// 初始化配置文件
func InitConfig() {

	v := viper.New()

	config := GetSystemConfig()

	if config {
		v.SetConfigFile("./config/pro_config.yaml")

	} else {
		v.SetConfigFile("./config/dev_config.yaml")
	}

	if err := v.ReadInConfig(); err != nil {
		zap.S().Fatal("读取文件失败")
	}

	if err := v.Unmarshal(&global.UserServerConfig); err != nil {
		zap.S().Fatal("配置文件获取失败")
	}
	log.Println("配置文件1", global.UserServerConfig)
	InitNacos()
	log.Println("配置文件2", global.UserServerConfig)

}

// 初始化nacos
func InitNacos() {
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.UserServerConfig.Nacos.Namespace, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      global.UserServerConfig.Nacos.ServerIp,
			ContextPath: "/nacos",
			Port:        global.UserServerConfig.Nacos.ServerPort,
			Scheme:      "http",
		},
	}
	client, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	content, _ := client.GetConfig(vo.ConfigParam{
		DataId: global.UserServerConfig.Nacos.DataId,
		Group:  global.UserServerConfig.Nacos.Group})

	json.Unmarshal([]byte(content), &global.UserServerConfig)
}
