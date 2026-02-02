package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/log"
	"github.com/zhilv666/linkchecker/pkg/response"
	"go.uber.org/zap"
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
	var req ListReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		log.Error("ListWithPageSize 参数不匹配", zap.Error(err))
		response.Fail(ctx, 400, "参数不匹配")
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 1
	} else if req.Size > 100 {
		req.Size = 100
	}

	linkRecord, count, err := l.service.ListWithPageSize(req.Page, req.Size, req.Keyword)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, map[string]any{
		"list":    linkRecord,
		"page":    req.Page,
		"size":    req.Size,
		"keyword": req.Keyword,
		"count":   count,
	})
}
