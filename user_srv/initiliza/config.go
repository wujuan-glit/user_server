package initiliza

import (
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

	log.Println("配置文件", global.UserServerConfig)
	//监听配置文件是否发生变化

}
