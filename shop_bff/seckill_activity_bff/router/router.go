package router

import (
	"github.com/gin-gonic/gin"
	"order_bff/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	InitShopCart(r)
	InitOrder(r)
	InitAlipay(r)
	r.GET("/health", api.Health)
	return r
}
