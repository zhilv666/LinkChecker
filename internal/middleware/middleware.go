package middleware

import "net/http"

// Next 定义了"下一步"执行的函数签名
// 它的底层类型其实就是 http.HandlerFunc，但改名后语义更清晰：
// "调用 next(w, r) 就是把控制权移交给下一个中间件"
type Next func(w http.ResponseWriter, r *http.Request)

// Core 是我们需要编写的核心业务逻辑函数
// 相比之前的 Handler，这里 next 参数使用了上面的 Next 类型
type Core func(w http.ResponseWriter, r *http.Request, next Next)

// Standard 是 Go 标准库/Chi 识别的标准中间件签名
// 定义这个别名是为了让返回值看起来不那么吓人
type Standard func(http.Handler) http.Handler

// Wrap 将我们的核心逻辑 (Core) 转换为标准中间件 (Standard)
func Wrap(c Core) Standard {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 技巧：由于 Next 和 http.HandlerFunc 的签名完全一致
			// 我们可以直接把 next.ServeHTTP 强转为 Next 类型传进去
			c(w, r, Next(next.ServeHTTP))
		})
	}
}
