package initiliza

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	zap.S().Info("配置文件", global.ServerConfig)

}
