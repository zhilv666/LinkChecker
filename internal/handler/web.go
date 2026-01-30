package handler

import (
	"html/template"
	"net/http"

	"github.com/zhilv666/linkchecker/internal/service"
)

type WebHandler struct {
	svc *service.LinkService
}

func NewWebHandler(svc *service.LinkService) *WebHandler {
	return &WebHandler{
		svc: svc,
	}
}

func (h *WebHandler) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "模板解析失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	list, _, err := h.svc.ListWithPageSize(1, 10)
	if err != nil {
		http.Error(w, "数据获取失败", http.StatusInternalServerError)
		return
	}
	data := map[string]any{
		"List": list,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "数据渲染失败: "+err.Error(), http.StatusInternalServerError)
	}
}
