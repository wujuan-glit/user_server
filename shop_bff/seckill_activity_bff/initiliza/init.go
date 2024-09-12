package initiliza

func Init() {
	InitConfig()
	InitLogger()
	InitRedis()

	InitServerConn()
	InitSentinel()
	InitRegisterConsul()
}
