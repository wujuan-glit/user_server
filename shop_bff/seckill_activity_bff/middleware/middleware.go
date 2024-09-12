package middleware

import (
	"encoding/json"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"net/http"
	"order_bff/model"
	"order_bff/server"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")

		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    999,
				"message": "token未传",
			})
			c.Abort()
			return

		}
		jwtToken, err := server.CheckJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    999,
				"message": "token失效",
			})
			c.Abort()
			return
		}
		var user model.User
		err = json.Unmarshal([]byte(jwtToken), &user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    999,
				"message": "结构体转换失败",
			})

			return
		}
		//把用户信息放入上下文中
		c.Set("user_id", user.ID)
		c.Next()

	}
}

// QpsLimit qbs限流

func QpsLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, b := sentinel.Entry("order_qps", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": -1,
				"msg":  "请求频繁,请稍后请求",
			})
			c.Abort()
			return
		}
	}
}

//func QpsLimit() gin.HandlerFunc {
//	return func(c *gin.Context) {
//
//		e, b := sentinel.Entry("order_limit", sentinel.WithTrafficType(base.Inbound))
//		if b != nil {
//			c.JSON(http.StatusTooManyRequests, gin.H{
//				"code":    -1,
//				"message": "限流了",
//			})
//			c.Abort()
//			return
//		}
//		// 确保在请求处理完毕后退出 Sentinel 入口点
//		defer func() {
//			if r := recover(); r != nil {
//				// 如果后续处理中发生 panic，则捕获并记录错误
//				sentinel.TraceError(e, fmt.Errorf("panic: %v", r))
//				// 根据需要设置响应或重新抛出 panic
//				c.AbortWithStatus(http.StatusInternalServerError)
//			}
//			e.Exit()
//		}()
//
//		// 继续执行后续中间件或处理器
//		c.Next()
//
//		// 如果后续处理中有错误，则记录到 Sentinel（注意：这里已经包含了 panic 的情况）
//		if c.Writer.Status() >= http.StatusInternalServerError {
//			sentinel.TraceError(e, errors.New("internal server error"))
//		}
//	}
//
//}
//
//func SentinelMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 尝试获取 Sentinel 的入口点
//
//		e, b := sentinel.Entry("order-qps", sentinel.WithTrafficType(base.Inbound))
//		if b != nil {
//			// 流量控制或熔断降级触发，拒绝请求
//			c.JSON(http.StatusTooManyRequests, gin.H{
//				"code":    -1,
//				"message": "请求过于频繁，请稍后再试",
//				"data":    nil,
//			})
//			c.Abort()
//			return
//		}
//
//		// 确保在请求处理完毕后退出 Sentinel 入口点
//		defer func() {
//			if r := recover(); r != nil {
//				// 如果后续处理中发生 panic，则捕获并记录错误
//				sentinel.TraceError(e, fmt.Errorf("panic: %v", r))
//				// 根据需要设置响应或重新抛出 panic
//				c.AbortWithStatus(http.StatusInternalServerError)
//			}
//			e.Exit()
//		}()
//
//		// 继续执行后续中间件或处理器
//		c.Next()
//
//		// 如果后续处理中有错误，则记录到 Sentinel（注意：这里已经包含了 panic 的情况）
//		if c.Writer.Status() >= http.StatusInternalServerError {
//			sentinel.TraceError(e, errors.New("internal server error"))
//		}
//	}
//}
