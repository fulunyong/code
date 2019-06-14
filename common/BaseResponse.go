package common

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//统一返回处理
func Render(w http.ResponseWriter, response BaseResponse) {
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(response)
	_, _ = w.Write(bytes)
}
