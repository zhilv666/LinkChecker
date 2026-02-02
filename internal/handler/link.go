package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/response"
)

type LinkHandler struct {
	service *service.LinkService
}

func NewLinkHandler(s *service.LinkService) *LinkHandler {
	return &LinkHandler{
		service: s,
	}
}

func (l *LinkHandler) CheckOne(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	url := query.Get("url")
	password := query.Get("password")

	if url == "" {
		response.Fail(ctx, 400, "链接不能为空")
		return
	}

	linkRecord, err := l.service.CheckOne(url, password)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, linkRecord)
}

func (l *LinkHandler) ListWithPageSize(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	pageStr := query.Get("page")
	sizeStr := query.Get("size")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		response.Fail(ctx, 400, err.Error())
		return
	}
	sizeInt, err := strconv.Atoi(sizeStr)
	if err != nil {
		response.Fail(ctx, 400, err.Error())
		return
	}
	if pageInt < 1 {
		pageInt = 1
	}
	if sizeInt < 1 {
		sizeInt = 1
	} else if sizeInt > 100 {
		sizeInt = 100
	}

	linkRecord, count, err := l.service.ListWithPageSize(pageInt, sizeInt)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, map[string]any{
		"list":  linkRecord,
		"page":  pageInt,
		"size":  sizeInt,
		"count": count,
	})
}
