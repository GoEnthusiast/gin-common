package responsex

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	CodeSuccess                   = NewCode(10000, "成功")
	CodeInvalidParams             = NewCode(10001, "参数错误")
	CodeServerError               = NewCode(10002, "服务器内部错误")
	CodeTooManyRequest            = NewCode(10003, "请求次数过多")
	CodeUnauthorizedSignError     = NewCode(10004, "sign错误")
	CodeUnauthorizedSignTimeout   = NewCode(10005, "sign过期")
	CodeUnauthorizedSignNotExist  = NewCode(10006, "sign缺失")
	CodeMessageError              = NewCode(10007, "Message Error")
	CodeUnauthorizedTokenError    = NewCode(10008, "token错误")
	CodeUnauthorizedTokenTimeout  = NewCode(10009, "token过期")
	CodeUnauthorizedTokenNotExist = NewCode(100010, "token缺失")
)

var codes = map[int]string{}

type E struct {
	code int
	msg  string
	data []string
}

func (e E) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s\n", e.Code(), e.Msg())
}

func (e *E) Code() int {
	return e.code
}

func (e *E) Msg() string {
	return e.msg
}

func (e *E) Details() []string {
	return e.data
}

func (e *E) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.Msg(), args...)
}

func (e *E) WithDetails(details ...string) *E {
	newError := *e
	newError.data = []string{}
	for _, d := range details {
		newError.data = append(newError.data, d)
	}
	return &newError
}

func (e *E) WithInterface(detail interface{}) *E {
	newError := *e
	detailByte, _ := json.Marshal(detail)
	newError.data = []string{}
	newError.data = append(newError.data, string(detailByte))
	return &newError
}

func (e *E) StatusCode() int {
	switch e.Code() {
	case CodeSuccess.Code():
		return http.StatusOK
	case CodeInvalidParams.Code():
		return http.StatusOK
	case CodeServerError.Code():
		return http.StatusOK
	case CodeTooManyRequest.Code():
		return http.StatusOK
	case CodeUnauthorizedSignError.Code():
		return http.StatusOK
	case CodeUnauthorizedSignTimeout.Code():
		return http.StatusOK
	case CodeUnauthorizedSignNotExist.Code():
		return http.StatusOK
	case CodeMessageError.Code():
		return http.StatusOK
	case CodeUnauthorizedTokenError.Code():
		return http.StatusOK
	case CodeUnauthorizedTokenTimeout.Code():
		return http.StatusOK
	case CodeUnauthorizedTokenNotExist.Code():
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}

func NewCode(code int, msg string) *E {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码「%d」重复。", code))
	}
	codes[code] = msg
	return &E{code: code, msg: msg}
}
