package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"user/user_buf/model"
	"user/user_buf/server"
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
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Next()

	}
}
