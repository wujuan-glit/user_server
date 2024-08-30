package router

import (
	"github.com/gin-gonic/gin"
	"order_bff/api/shop_cart"
	"order_bff/middleware"
)

func InitShopCart(r *gin.Engine) {
	v1 := r.Group("o/v1/cart").Use(middleware.Auth(), middleware.SentinelMiddleware())
	v1.GET("/list", shop_cart.CartList)
	v1.POST("/add", shop_cart.CreateCart)
	v1.POST("/update", shop_cart.UpdateCart)
	v1.POST("/delete", shop_cart.DeleteCart)
}
