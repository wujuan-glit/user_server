package router

import (
	"github.com/gin-gonic/gin"
	"order_bff/api/alipay"
	"order_bff/middleware"
)

func InitAlipay(r *gin.Engine) {
	v1 := r.Group("o/v1/alipay").Use(middleware.SentinelMiddleware())
	v1.POST("/generate", alipay.GenerateAlipayUrl)
	v1.GET("/return", alipay.ReturnUrl)
	v1.Any("/notify", alipay.NotifyUrl)

}
