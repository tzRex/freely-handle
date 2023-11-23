package corecode

type baseReq struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	// 一切正常
	CodeSuccess = 1000
	// 参数类型错误
	CodeFail = 1001
	// 数据未通过验证
	CodeValidateFail = 1002
)

var (
	CodeRdbTokenCurrent = "-current"
	CodeRdbTokenBackup  = "-backup"
)

func ReqOk(data interface{}) *baseReq {
	return &baseReq{Code: CodeSuccess, Msg: "操作成功", Data: data}
}

func ReqFail(msg string, tip ...string) *baseReq {
	return &baseReq{Code: CodeValidateFail, Msg: msg}
}

func ReqBad(msg string, tip ...string) *baseReq {
	return &baseReq{Code: CodeFail, Msg: msg}
}
