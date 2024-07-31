package router

import (
	"github.com/gin-gonic/gin"
	"user/user_buf/api"
	"user/user_buf/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("u/v1/user")
	{
		v1.POST("/register", api.RegisterUser)
		v1.GET("/list", api.GetUserList)
		v1.PUT("/update", middleware.Auth(), api.UpdateUserInfo)
		v1.POST("/login", api.Login)
		v1.POST("/refresh", api.ReFresh)
	}
	v2 := r.Group("captcha")
	{
		v2.GET("", api.GetCaptcha)
	}
	r.GET("/health", api.Health)
	return r
}
