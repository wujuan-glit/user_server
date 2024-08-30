package main

import (
	"fmt"
	"go.uber.org/zap"
	"user/user_buf/global"
	"user/user_buf/initiliza"

	"user/user_buf/router"
)

func main() {

	initiliza.Init()

	r := router.InitRouter()

	err := r.Run(fmt.Sprintf("%s:%d", global.ServerConfig.PubIp, global.ServerConfig.Port))
	if err != nil {
		zap.S().Info("启动失败", err)
	}

	//等待中断以后优雅的关闭服务器（设置5s的超时时间）

	/*quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	_, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancelFunc()*/

}
