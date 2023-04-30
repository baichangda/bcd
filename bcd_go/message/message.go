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
	marshal, err := json.Marshal(Succeed_msg(msg))
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
	_, err = ctx.Writer.Write(marshal)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
}

func ResponseSucceed_data(data interface{}, ctx *gin.Context) {
	bs, err := json.Marshal(Succeed_data(data))
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
	_, err = ctx.Writer.Write(bs)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}

}

func ResponseFailed_msg(code int, msg string, ctx *gin.Context) {
	marshal, err := json.Marshal(Failed(code, msg))
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
	_, err = ctx.Writer.Write(marshal)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
}

func ResponseFailed_err(err error, ctx *gin.Context) {
	util.Log.Errorf("%+v", err)
	failed := Failed(1, err.Error())
	marshal, err := json.Marshal(failed)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
	_, err = ctx.Writer.Write(marshal)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
}
