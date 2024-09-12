package router

import (
	"github.com/gin-gonic/gin"
	"order_bff/api/order"
	"order_bff/middleware"
)

func InitOrder(r *gin.Engine) {
	v1 := r.Group("o/v1/order").Use(middleware.Auth())
	v1.GET("/list", order.ListOrder)
	v1.POST("/add", order.CreateOrder)
	v1.GET("/update", order.UpdateOrder)

}
