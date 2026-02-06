package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhilv666/linkchecker/internal/dto"
	"github.com/zhilv666/linkchecker/internal/service"
	"github.com/zhilv666/linkchecker/pkg/response"
)

type LinkHandler struct {
	svc *service.LinkService
}

func NewLinkHandler(svc *service.LinkService) *LinkHandler {
	return &LinkHandler{svc: svc}
}

func (h *LinkHandler) CheckOne(ctx *gin.Context) {
	var req dto.CheckReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, 400, "参数错误")
		return
	}

	if req.URL == "" {
		response.Fail(ctx, 400, "链接不能为空")
		return
	}
	data, err := h.svc.CheckAndSave(ctx.Request.Context(), req.URL, req.Password)

	if err != nil {
		response.Fail(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, data)
}

func (h *LinkHandler) ListWithPageSize(ctx *gin.Context) {
	var req dto.ListReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, 400, "参数错误")
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

	linkRecord, count, err := h.svc.GetLinkList(ctx.Request.Context(), req.Page, req.Size, req.Keyword)
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

func (h *LinkHandler) Report(ctx *gin.Context) {
	var req dto.ReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, 400, "参数错误")
		return
	}

	err := h.svc.SaveResult(ctx, &req)
	if err != nil {
		response.Fail(ctx, 500, err.Error())
		return
	}
	response.Success(ctx, nil)
}
