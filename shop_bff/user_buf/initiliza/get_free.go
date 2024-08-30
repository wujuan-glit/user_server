package initiliza

import (
	"go.uber.org/zap"
	"user/user_buf/global"

	"net"
)

func InitFreePort() {

	if global.ServerConfig.Port == 0 {
		port, err := GetFreePort()
		if err != nil {
			zap.S().Panic(err)
		}
		//重新赋值
		global.ServerConfig.Port = port
	}
}
func GetFreePort() (int, error) {
	//函数用于解析TCP地址
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		zap.S().Panic("解析地址错误", err)
		return 0, err
	}
	//在指定的地址上监听TCP网络连接
	tcp, err := net.ListenTCP("tcp", addr)
	if err != nil {
		zap.S().Panic("连接错误", err)
		return 0, err
	}

	return tcp.Addr().(*net.TCPAddr).Port, nil
}
