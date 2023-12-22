package errorx

type ResCode int64

const (
	CodeSuccess ResCode = 200
)

// standardization error
const (
	CodeHasNoRoute     ResCode = 404
	CodeNotAllowMethod ResCode = 408
)

// usually error
const (
	CodeServerBusy ResCode = 500 + iota
	CodeInvalidParam
	CodeErrAuth
	CodeTokenExpire
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:        "成功",
	CodeHasNoRoute:     "路由不存在",
	CodeServerBusy:     "服务繁忙",
	CodeInvalidParam:   "请求参数错误",
	CodeNotAllowMethod: "方法不允许",
	CodeErrAuth:        "权限不足",
	CodeTokenExpire:    "token过期",
}

func (r ResCode) Msg() string {
	msg, ok := codeMsgMap[r]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}

func (r ResCode) Error() string {
	return r.Msg()
}
