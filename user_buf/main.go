package main

import (
	"go.uber.org/zap"

	"user/user_buf/initiliza"

	"user/user_buf/router"
)

func main() {

	initiliza.InitServerConn()
	initiliza.InitTranslation("zh")
	initiliza.InitRegisterValidator()

	r := router.InitRouter()
	err := r.Run("127.0.0.1:8888")
	if err != nil {
		zap.S().Info("启动失败", err)
	}
}
