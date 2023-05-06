package message

import (
	"bcd_go/util"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type JsonMessage struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (msg *JsonMessage) Response(ctx *gin.Context) {
	marshal, err := json.Marshal(msg)
	if err != nil {
		util.Log.Errorf("%+v", err)
	}
	_, err = ctx.Writer.Write(marshal)
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

func ResponseSucceed_msg(msg string, ctx *gin.Context) {
	Succeed_msg(msg).Response(ctx)
}

func ResponseSucceed_data(data interface{}, ctx *gin.Context) {
	Succeed_data(data).Response(ctx)
}

func FromGinError(err *gin.Error) *JsonMessage {
	meta, ok := err.Meta.(*ErrorMeta)
	if ok {
		return &JsonMessage{
			Code:    meta.Code,
			Message: err.Error(),
			Data:    meta.Data,
		}
	} else {
		return &JsonMessage{
			Code:    1,
			Message: err.Error(),
		}
	}
}

func GinError_msg_code(ctx *gin.Context, msg string, code int) {
	_ = ctx.Error(&gin.Error{
		Err:  errors.New(msg),
		Type: gin.ErrorTypeAny,
		Meta: ErrorMeta{
			Code: code,
			Data: nil,
		},
	})
}

func GinError_msg(ctx *gin.Context, msg string) {
	GinError_msg_code(ctx, msg, 1)
}

func GinError_err(ctx *gin.Context, err error) {
	_ = ctx.Error(errors.WithStack(err))
}

type ErrorMeta struct {
	Code int
	Data any
}
