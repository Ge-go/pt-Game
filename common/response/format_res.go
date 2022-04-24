package response

type StandardRes struct {
	Ret  int         `json:"ret"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func GenSuccessRes(data interface{}) *StandardRes {
	return &StandardRes{Ret: 0, Msg: "Success", Data: data}
}

func GenFailedRes(errCode int, errMsg string) *StandardRes {
	return &StandardRes{Ret: errCode, Msg: errMsg, Data: new(map[string]interface{})}
}
