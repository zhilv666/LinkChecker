package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/pkg/response"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit // 每秒加入的令牌数
	b   int        // 桶大小
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	if limiter, exists := i.ips[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter
	return limiter
}

func RateMiddleware(limit *IPRateLimiter) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		limiter := limit.GetLimiter(ip)

		if !limiter.Allow() {
			ctx.JSON(http.StatusOK, response.Response{
				Code: 429,
				Msg:  "请求过于频繁，请稍后再试",
				Data: nil,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
