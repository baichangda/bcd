package ocr

import (
	"bcd_go/baidu"
	"bcd_go/message"
	"bcd_go/tencent"
	"bcd_go/util"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	v20181119 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
	"strconv"
	"strings"
	"time"
)

func Route(engine *gin.Engine) {
	engine.POST("/api/ocr/image/baiduAI", process_image_baiduAI)
	engine.POST("/api/ocr/image/baiduFanyi", process_image_baiduFanyi)
	engine.POST("/api/ocr/table/baiduAI", process_table_baiduAI)
	engine.POST("/api/ocr/table/tencentAI", process_table_tencetAI)
}

func process_image_baiduAI(ctx *gin.Context) {
	//cfg.Log.Debugf("start process_image_baiduAI")
	buffer := bytes.Buffer{}
	_, err := buffer.ReadFrom(ctx.Request.Body)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	split := strings.Split(buffer.String(), ",")
	json, err := baidu.OcrAccurate(split[1], "", "", "", split[0], "", "", "")
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	words_result := json.Get("words_result")
	if words_result.Exists() {
		sb := strings.Builder{}
		for _, cur := range words_result.Array() {
			sb.WriteString(cur.Get("words").Str)
			sb.WriteString("\n")
		}
		message.ResponseSucceed_data(sb.String(), ctx)
	} else {
		message.GinError_msg(ctx, fmt.Sprintf("失败、错误信息:\n%s", json.Raw))
	}
}

func process_image_baiduFanyi(ctx *gin.Context) {
	buffer := bytes.Buffer{}
	_, err := buffer.ReadFrom(ctx.Request.Body)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	split := strings.Split(buffer.String(), ",")

	bs, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	json, err := baidu.OcrFanyi(split[0], &bs)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	errno := json.Get("errno")
	if errno.Exists() && errno.Int() == 0 {
		arr := json.Get("data").Get("src")
		sb := strings.Builder{}
		for _, cur := range arr.Array() {
			sb.WriteString(strconv.Quote(cur.Str))
			sb.WriteString("\n")
		}
		message.ResponseSucceed_data(sb.String(), ctx)
	} else {
		util.Log.Warn(json.Raw)
		message.GinError_msg(ctx, strconv.Quote(json.Get("errmsg").Str))
	}
}

const timeout = 15 * time.Second

func process_table_baiduAI(ctx *gin.Context) {
	buffer := bytes.Buffer{}
	_, err := buffer.ReadFrom(ctx.Request.Body)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	split := strings.Split(buffer.String(), ",")

	json, err := baidu.OcrFormAsync(split[1], "", "excel")
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	//cfg.Log.Debugf("json:\n%s", json.Raw)
	error_code := json.Get("error_code")
	if error_code.Exists() {
		message.GinError_msg(ctx, error_code.Raw)
	} else {
		request_id := json.Get("result").Array()[0].Get("request_id").Str
		end := time.Now().UnixMilli() + timeout.Milliseconds()
		for {
			resultJson, err := baidu.OcrFormAsyncResult(request_id, "excel")
			if err != nil {
				_ = ctx.Error(errors.WithStack(err))
				return
			}
			resultError_code := resultJson.Get("error_code")
			//cfg.Log.Debugf("resultJson:\n%s", resultJson.Raw)
			if resultError_code.Exists() {
				message.GinError_msg(ctx, resultError_code.Raw)
				return
			} else {
				retCode := resultJson.Get("result").Get("ret_code").Int()
				if retCode == 3 {
					result_data := resultJson.Get("result").Get("result_data")
					if result_data.Exists() {
						message.ResponseSucceed_data(result_data.Str, ctx)
					} else {
						message.GinError_msg(ctx, fmt.Sprintf("执行失败\n%s", resultJson.Raw))
					}
					break
				} else {
					if time.Now().UnixMilli() >= end {
						message.GinError_msg(ctx, "获取结果超时")
						break
					} else {
						time.Sleep(2 * time.Second)
					}
				}

			}

		}
	}
}

func process_table_tencetAI(ctx *gin.Context) {
	buffer := bytes.Buffer{}
	_, err := buffer.ReadFrom(ctx.Request.Body)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	split := strings.Split(buffer.String(), ",")
	request := v20181119.NewRecognizeTableOCRRequest()
	request.ImageBase64 = &split[1]
	request.TableLanguage = &split[0]
	response, err := tencent.ClientOcrTable.RecognizeTableOCR(request)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	dataBase64 := response.Response.Data
	bs, err := base64.StdEncoding.DecodeString(*dataBase64)
	if err != nil {
		_ = ctx.Error(errors.WithStack(err))
		return
	}
	ctx.Set("content-type", "application/octet-stream")
	ctx.Set("content-disposition", "attachment;filename=table.xlsx")
	ctx.Writer.Write(bs)
}
