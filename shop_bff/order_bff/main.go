package main

import (
	"fmt"
	"go.uber.org/zap"
	"order_bff/api/alipay"
	"order_bff/initiliza"
	"order_bff/router"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initiliza.Init()
	r := router.InitRouter()
	//go alipay.AlipayConsumer()

	alipay.RocketMqConsumer()
	go func() {
		err := r.Run(fmt.Sprintf("%s:%d", "10.3.189.2", 8891))

		if err != nil {
			zap.S().Panic("启动失败")
		}
	}()
	//优雅退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	initiliza.Deregister()

}
