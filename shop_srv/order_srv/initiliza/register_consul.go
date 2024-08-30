package initiliza

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"order_srv/global"
)

func InitRegisterConsul() {

	//创建一个Consul客户端配置  这里的ip和端口是用来连接consul的  consul默认的是8500
	cfg := api.DefaultConfig()
	// 设置Consul客户端的地址，从全局配置中获取Consul的主机和端口
	cfg.Address = fmt.Sprintf("%s:%d", global.UserServerConfig.Consul.Host, global.UserServerConfig.Consul.Port)
	//使用配置创建Consul客户端实例
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 创建一个服务检查对象，用于监控服务状态

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.UserServerConfig.Ip, global.UserServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	// 创建一个服务注册对象，用于向Consul注册服务信息

	res := api.AgentServiceRegistration{
		Name:    global.UserServerConfig.Name,
		ID:      global.ServerID,
		Port:    global.UserServerConfig.Port,
		Tags:    global.UserServerConfig.Consul.Tags,
		Address: global.UserServerConfig.Ip,
		Check:   check,
	}
	// 使用Consul客户端注册服务
	err = client.Agent().ServiceRegister(&res)
	if err != nil {
		zap.S().Panic(err)
	}

}

// Deregister 注销consul服务
func Deregister() {
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
}
