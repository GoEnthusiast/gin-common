package responsex

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type result struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Body  interface{} `json:"body"`
	Error []string    `json:"error"`
}

type StreamLogicResult struct {
	Data  string
	Error error
}

// WriteJson
/*
	返回Json响应
*/
func WriteJson(ginCtx *gin.Context, resp interface{}, err error) {
	if err == nil {
		ginCtx.JSON(CodeSuccess.StatusCode(), &result{
			Code:  CodeSuccess.Code(),
			Msg:   CodeSuccess.Msg(),
			Body:  resp,
			Error: []string{},
		})
	} else {
		e, ok := err.(*E)
		if !ok {
			zap.L().Error(fmt.Sprintf("响应的error不是<*status.Error>类型, 请求的url: %s", ginCtx.Request.URL.Path))
			er := CodeServerError
			ginCtx.JSON(er.StatusCode(), &result{
				Code:  er.Code(),
				Msg:   er.Msg(),
				Body:  nil,
				Error: er.Details(),
			})
			return
		}
		ginCtx.JSON(e.StatusCode(), &result{
			Code:  e.Code(),
			Msg:   e.Msg(),
			Body:  nil,
			Error: e.Details(),
		})
	}
}

// WriteStream
/*
	流式响应
*/
func WriteStream(ginCtx *gin.Context, f func(ch *chan []byte, args ...any), args ...any) error {
	flusher, ok := ginCtx.Writer.(http.Flusher)
	if !ok {
		return errors.New("Streaming not supported")
	}
	ginCtx.Writer.Header().Set("Content-Type", "text/plain")
	ginCtx.Writer.Header().Set("Cache-Control", "no-cache")
	ginCtx.Writer.Header().Set("Connection", "keep-alive")
	var chanResult = make(chan []byte)
	go f(&chanResult, args...)
	for {
		res, ok := <-chanResult
		if !ok {
			break
		}
		ginCtx.Writer.Write(res)
		flusher.Flush()
	}
	return nil
}
