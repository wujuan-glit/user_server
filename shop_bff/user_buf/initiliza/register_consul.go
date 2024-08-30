package initiliza

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"user/user_buf/global"
)

func InitRegisterConsul() {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.Consul.Host, global.ServerConfig.Consul.Port)

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
