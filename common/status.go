package common

//公共常量使用
const (
	ResponseOK    = 0  //返回正常  其他为非正常
	ResponseError = -1 //返回失败
)

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
