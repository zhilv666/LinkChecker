package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/pkg/response"
)

func TokenMiddleware(_token string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("token")
		if token == "" {
			response.Fail(ctx, 400, "参数错误")
			ctx.Abort()
			return
		}

		if token != _token {
			response.Fail(ctx, 400, "参数错误")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
