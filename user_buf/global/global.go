package global

import (
	"context"
	ut "github.com/go-playground/universal-translator"
	"user/proto"
	"user/user_buf/config"
)

var (
	ServerConn proto.UserClient
	Ctx        context.Context

	ServerConfig config.UserServer
	RootPath     string
	// 定义一个全局翻译器T
	Trans ut.Translator
)
