package message

import (
	"bcd_go/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type JsonMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (msg *JsonMessage) ToBytes() []byte {
	marshal, err := json.Marshal(msg)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
	return marshal
}

func (msg *JsonMessage) Response(ctx *gin.Context) {
	marshal := msg.ToBytes()
	_, err := ctx.Writer.Write(marshal)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
}

func Succeed_msg(msg string) *JsonMessage {
	return &JsonMessage{
		Code:    0,
		Message: msg,
	}
}

func Succeed_data(data interface{}) *JsonMessage {
	return &JsonMessage{
		Code: 0,
		Data: data,
	}
}

func Failed(code int, message string) *JsonMessage {
	return &JsonMessage{
		Code:    code,
		Message: message,
	}
}

func ResponseSucceed_msg(msg string, ctx *gin.Context) {
	Succeed_msg(msg).Response(ctx)
}

func ResponseSucceed_data(data interface{}, ctx *gin.Context) {
	Succeed_data(data).Response(ctx)
}

func ResponseFailed_msg(code int, msg string, ctx *gin.Context) {
	Failed(code, msg).Response(ctx)
}

func ResponseFailed_err(err error, ctx *gin.Context) {
	Failed(1, err.Error()).Response(ctx)
}

type MyErr struct {
	Msg  string
	Code int
}

func (m *MyErr) Error() string {
	return m.Msg
}

func (m *MyErr) ToJsonMessage() *JsonMessage {
	return &JsonMessage{
		Code:    m.Code,
		Message: m.Msg,
	}
}

func NewMyError(msg string, code int) *MyErr {
	return &MyErr{
		Msg:  msg,
		Code: code,
	}
}
