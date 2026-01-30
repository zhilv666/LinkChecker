package middleware

import (
	"net/http"
	"time"

	"github.com/zhilv666/linkchecker/pkg/log"
	"go.uber.org/zap"
)

// 1. 定义包装器 (修正拼写 logging)
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// 2. 重写 WriteHeader (修正拼写 WriterHeader -> WriteHeader)
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// 3. 重写 Write
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

// 4. Logger 中间件 (使用 Wrap 简化)
func LoggerMiddleware() Standard {
	return Wrap(func(w http.ResponseWriter, r *http.Request, next Next) {
		start := time.Now()

		// 初始化包装器
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // 默认 200
		}

		next(lrw, r)

		// 计算耗时
		duration := time.Since(start)

		// 组装日志字段
		fields := []zap.Field{
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("query", r.URL.RawQuery),
			zap.String("ip", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.Int("status", lrw.statusCode),
			zap.Int("size", lrw.size),
			zap.Duration("cost", duration),
		}

		// 根据状态码分级打印
		if lrw.statusCode >= 500 {
			log.Error("Server Error", fields...)
		} else if lrw.statusCode >= 400 {
			log.Warn("Client Error", fields...)
		} else {
			log.Info("Access Log", fields...)
		}
	})
}
