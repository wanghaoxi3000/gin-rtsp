package serializer

import (
	"ginrtsp/util"

	"github.com/gin-gonic/gin"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// Err 通用错误处理
func Err(errCode int, msg string, err error) *Response {
	res := Response{
		Code: errCode,
		Msg:  msg,
	}

	if err != nil {
		util.Log().Error(err.Error())
		// 生产环境隐藏底层报错
		if gin.Mode() != gin.ReleaseMode {
			res.Error = err.Error()
		}
	}
	return &res
}
