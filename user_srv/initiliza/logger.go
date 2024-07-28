package initiliza

import "go.uber.org/zap"

func NewLogger() (*zap.Logger, error) {
	//读取环境变量中的值
	config := GetSystemConfig()

	if config {
		//生产环境
		cfg := zap.NewProductionConfig()

		//判断logger目录是否存在  不存在进行创建
		cfg.OutputPaths = []string{
			"./logger/user_log.log",
		}
		return cfg.Build()
	} else {
		//开发环境
		cfg, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}
}

// 初始化日志文件
func InitLogger() {
	logger, err := NewLogger()
	if err != nil {
		zap.S().Fatal("日志失败", err)
	}

	zap.ReplaceGlobals(logger)
}
