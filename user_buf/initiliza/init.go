package initiliza

func Init() {
	InitConfig()
	InitNacos()
	GetNacosConfig()
	InitServerConn()
	InitTranslation("zh")
	InitRegisterValidator()
	InitFreePort()
	InitRegisterConsul()

}
