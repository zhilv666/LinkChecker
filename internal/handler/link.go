package handler

import (
	"net/http"
	"strconv"

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

func (l *LinkHandler) CheckOne(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	url := query.Get("url")
	password := query.Get("password")

	if url == "" {
		response.Fail(w, 400, "链接不能为空")
		return
	}

	linkRecord, err := l.service.CheckOne(url, password)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, linkRecord)
}

func (l *LinkHandler) ListWithPageSize(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	pageStr := query.Get("page")
	sizeStr := query.Get("size")
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		response.Fail(w, 400, err.Error())
		return
	}
	sizeInt, err := strconv.Atoi(sizeStr)
	if err != nil {
		response.Fail(w, 400, err.Error())
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
		response.Error(w, err)
		return
	}
	response.Success(w, map[string]any{
		"list":  linkRecord,
		"page":  pageInt,
		"size":  sizeInt,
		"count": count,
	})
}
