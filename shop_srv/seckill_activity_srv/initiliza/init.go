package initiliza

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"userop_srv/global"
)

// 获取系统的环境变量
func GetSystemConfig() bool {
	viper.AutomaticEnv()

	getString := viper.GetString("BOOK")
	zap.S().Info("getString ", getString)

	if getString == "dev" {
		return false
	} else if getString == "prod" {
		return true
	}

	return false
}
func Init() {

	InitConfig()

	global.ServerID = fmt.Sprintf("%s", uuid.NewV4()) //唯一的ID

	InitFreePort()

	InitLogger()

	InitMysql()

	InitRegisterConsul()

	InitConn()
}
