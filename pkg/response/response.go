package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(w http.ResponseWriter, data any) {
	Json(w, http.StatusOK, 0, "success", data)
}

func Fail(w http.ResponseWriter, code int, msg string) {
	Json(w, http.StatusOK, code, msg, nil)
}

func Error(w http.ResponseWriter, err error) {
	Json(w, http.StatusInternalServerError, 500, err.Error(), nil)
}

func Json(w http.ResponseWriter, httpStatus int, code int, msg string, data any) {
	w.Header().Set("Content-Type", "applition/json")
	w.WriteHeader(httpStatus)
	resp := Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
