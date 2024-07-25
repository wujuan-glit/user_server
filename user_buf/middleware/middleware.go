package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
	"user/user_buf/global"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		md := metadata.New(map[string]string{
			"token": "123",
			"hgi":   "fefe",
		})
		global.Ctx = metadata.NewOutgoingContext(context.Background(), md)

	}
}
