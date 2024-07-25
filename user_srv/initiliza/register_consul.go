package initiliza

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"user/user_srv/global"
)

func InitRegisterConsul() {
	cfg := api.DefaultConfig()

	cfg.Address = fmt.Sprintf("%s:%d", global.UserServerConfig.Consul.Host, global.UserServerConfig.Consul.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象

	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.UserServerConfig.Ip, global.UserServerConfig.Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.UserServerConfig.Name
	registration.ID = global.ServerID
	registration.Port = global.UserServerConfig.Port
	registration.Tags = global.UserServerConfig.Consul.Tags
	registration.Address = global.UserServerConfig.Ip
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		zap.S().Panic(err)
	}

}
