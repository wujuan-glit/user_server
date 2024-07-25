package initiliza

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"user/proto"
	"user/user_buf/global"
)

/*
func InitClient() {

	conn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	global.ServerConn = proto.NewUserClient(conn)

}
*/
// 把拨号连接改成通过consul读取
func InitServerConn() {
	//进行consul的连接
	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"

	//实例化consul的服务端
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panic("实例化失败", err)
	}
	//进行服务过滤
	filter, err := client.Agent().ServicesWithFilter("Service == user_srv")
	if err != nil {
		zap.S().Panic(err)
		return
	}

	var address string
	var port int

	for _, service := range filter {
		address = service.Address
		port = service.Port
	}

	//进行拨号
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", address, port), grpc.WithInsecure())
	if err != nil {
		return
	}
	global.ServerConn = proto.NewUserClient(conn)
}
