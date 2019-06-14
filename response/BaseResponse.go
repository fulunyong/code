package response

import (
	"encoding/json"
	"github.com/fulunyong/code/common"
	"net/http"
)

//统一返回处理
func Render(w http.ResponseWriter, response BaseResponse) {
	w.WriteHeader(http.StatusOK)
	bytes, _ := json.Marshal(response)
	_, _ = w.Write(bytes)
}
