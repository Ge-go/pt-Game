package response

import "github.com/kataras/iris/v12"

type Response struct {
	Ret  int         `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Send(c iris.Context, err error, data interface{}) {
	statusCode, serviceCode, message := DecodeErr(c, err)
	c.StopWithJSON(statusCode, Response{
		Ret:  serviceCode,
		Msg:  message,
		Data: data,
	})
}

// Res 主要用于生成swagger文档中统一的响应类型（swagger 无法解析interface类型，所以无法使用Response struct）
type Res struct {
	Ret  int                    `json:"ret"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}
