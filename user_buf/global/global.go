package global

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	uuid "github.com/satori/go.uuid"
	"user/proto"
	"user/user_buf/config"
)

var (
	ServerConn   proto.UserClient
	NacosClient  config_client.IConfigClient
	ServerConfig config.UserServer
	// 定义一个全局翻译器T
	Trans    ut.Translator
	ServerId string = fmt.Sprintf("%s", uuid.NewV4())
)
