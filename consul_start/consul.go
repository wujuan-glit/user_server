package main

import (
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

func RegisterConsul() {
	cfg := api.DefaultConfig()

	cfg.Address = "127.0.0.1:8500"
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           "10.3.189.2:8081",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}
	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = "user_srv"
	registration.ID = "user_srv2"
	registration.Port = 1000
	registration.Tags = []string{"wujuan", "user_srv"}
	registration.Address = "127.0.0.1"
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Panic(err)
	}
}
func main() {
	RegisterConsul()
}
