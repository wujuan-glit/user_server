package api

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

// 生成图形验证码
func GetCaptcha(c *gin.Context) {
	// 创建driver 高 宽 长度 干扰线数量 嘈点
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)

	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.DefaultMemStore)
	id, url, answer, err := captcha.Generate()
	if err != nil {
		ReturnErrorJson(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":       0,
		"message":    "success",
		"captcha_id": id,
		"url":        url,
		"answer":     answer,
	})
}
