package initiliza

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	useropPb "github.com/wujuan-glit/shop/userop"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"userop_srv/global"
	"userop_srv/handler"
)

// 初始化grpc连接
func InitConn() {

	g := grpc.NewServer()

	s := handler.UseropServer{}

	useropPb.RegisterAddressServer(g, &s)
	useropPb.RegisterMessageServer(g, &s)
	useropPb.RegisterUserFavServer(g, &s)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", global.UserServerConfig.Ip, global.UserServerConfig.Port))

	grpc_health_v1.RegisterHealthServer(g, health.NewServer())

	if err != nil {
		return
	}

	// 启动 gRPC 服务器并开始监听。
	go func() {
		err = g.Serve(listen)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//如果用户微服务注销了，就应该服务注册发现也注销了
	// 注销服务
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	//接收终止信号

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.UserServerConfig.Consul.Host, global.UserServerConfig.Consul.Port)

	//实例化consul客户端
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panic(err)
	}
	//注销服务
	err = client.Agent().ServiceDeregister(global.ServerID)
	if err != nil {
		zap.S().Info("服务注销失败")
		return
	}

	zap.S().Info("服务注销成功")

}
