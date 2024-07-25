package api

import (
	"net/http"
	"strconv"
	"strings"
	"user/user_buf/global"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ReturnErrorJson(c *gin.Context, err error) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)

	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 999,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	// validator.ValidationErrors类型错误则进行翻译
	//并使用removeTopStruct函数去除字段名中的结构体名称标识
	c.JSON(http.StatusBadRequest, gin.H{
		"code": 999,
		"msg":  removeTopStruct(errs.Translate(global.Trans)),
		"data": nil,
	})
	return
}

// removeTopStruct函数去除字段名中的结构体名称标识
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func HandleGrpcErrorToHttp(c *gin.Context, err error) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"code": 888,
					"data": nil,
					"msg":  e.Message(),
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusNotFound, gin.H{
					"code": 888,
					"data": nil,
					"msg":  e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 888,
					"data": nil,
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 888,
					"data": nil,
					"msg":  "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 888,
					"data": nil,
					"msg":  "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 888,
					"data": nil,
					"msg":  e.Code(),
				})
			}
			return
		}
	}
}
func GetPage(c *gin.Context) (int, int, int) {
	page := c.Query("page")
	limit := c.Query("limit")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	if pageInt == 0 {
		pageInt = 1
	}

	if limitInt == 0 {
		limitInt = 10
	}

	if limitInt > 100 {
		limitInt = 100
	}

	offset := (pageInt - 1) * limitInt
	return pageInt, limitInt, offset
}