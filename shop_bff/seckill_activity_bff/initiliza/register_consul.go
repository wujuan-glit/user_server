package initiliza

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"order_bff/global"
)

//订单 order_web

func InitRegisterConsul() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("127.0.0.1:8500")

	client, err := api.NewClient(cfg)
	if err != nil {
		return
	}

	check := &api.AgentServiceCheck{
		Interval:                       "5s",
		Timeout:                        "5s",
		HTTP:                           fmt.Sprintf("http://%s:%d/health", global.ServerConfig.PubIp, global.ServerConfig.Port),
		DeregisterCriticalServiceAfter: "10s",
	}

	registration := api.AgentServiceRegistration{
		ID:      global.ServerId,
		Name:    global.ServerConfig.Name,
		Tags:    global.ServerConfig.Consul.Tags,
		Port:    global.ServerConfig.Port,
		Address: global.ServerConfig.PubIp,
		Check:   check,
	}

	err = client.Agent().ServiceRegister(&registration)
	if err != nil {
		zap.S().Panic("consul服务注册失败err", err)
	}
}

// Deregister 注销consul服务
func Deregister() {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)

	//实例化consul客户端
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panic(err)
	}
	//注销服务
	err = client.Agent().ServiceDeregister(global.ServerId)
	if err != nil {
		zap.S().Info("服务注销失败")
		return
	}

}
