package router

import (
	"github.com/gin-gonic/gin"
	"user/user_buf/api"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("user")
	{
		v1.POST("/register", api.RegisterUser)
		v1.GET("/list", api.GetUserList)
		v1.PUT("/update", api.UpdateUserInfo)
		v1.POST("/login", api.Login)
	}
	v2 := r.Group("captcha")
	{
		v2.GET("", api.GetCaptcha)
	}
	return r
}
