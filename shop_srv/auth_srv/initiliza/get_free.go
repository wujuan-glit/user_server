package initiliza

import (
	"go.uber.org/zap"
	"net"
	"user_srv/global"
)

func InitFreePort() {

	if global.UserServerConfig.Port == 0 {
		port, err := GetFreePort()
		if err != nil {
			zap.S().Panic(err)
		}
		global.UserServerConfig.Port = port

	}

}

func GetFreePort() (int, error) {
	//函数用于解析TCP地址
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		zap.S().Panic("err", err)
		return 0, err
	}
	//函数用于在指定的地址上监听TCP网络连接
	l, err := net.ListenTCP("tcp", addr)

	if err != nil {
		return 0, err
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port, nil
}
