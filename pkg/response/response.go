package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Fail(ctx *gin.Context, code int, msg string) {
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, Response{
		Code: 500,
		Msg:  err.Error(),
		Data: nil,
	})
}

func Json(ctx *gin.Context, httpStatus int, code int, msg string, data any) {
	ctx.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
