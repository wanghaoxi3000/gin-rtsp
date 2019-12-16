package api

import (
	"encoding/json"
	"fmt"
	"ginrtsp/serializer"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

// errorRequest 请求数据错误处理
func errorRequest(err error) *serializer.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			return serializer.Err(
				400,
				fmt.Sprintf("%s %s", e.Field, e.Tag),
				err,
			)
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Err(400, "JSON类型不匹配", err)
	}

	return serializer.Err(400, "参数错误", err)
}

// Ping 状态检查
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}
