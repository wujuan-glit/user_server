package initiliza

func Init() {
	InitConfig()
	InitServerConn()
	InitTranslation("zh")
	InitRegisterValidator()
	InitFreePort()

}
