package initiliza

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化sentinel 异常: %v", err)
	}
	//配置限流规则
	//这种配置应该从nacos中读取
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "order-qps",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject, //直接拒绝
			Threshold:              1,           //1秒钟允许访问的商品列表为1次
			StatIntervalInMs:       1000,
		},
	})

}
