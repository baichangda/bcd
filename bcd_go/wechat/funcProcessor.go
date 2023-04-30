package wechat

import (
	"bcd_go/baidu"
	"bcd_go/util"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var allFuncDesc = []string{
	//"根据图片识别车型",
	//"根据图片识别车辆损伤",
	//"人像动漫化",
}
var allFuncTip = []string{
	//"请上传图片",
	//"请上传图片",
	//"请上传图片",
}
var allFunc = []func(inMessage InMessage) (*OutMessage, error){
	//processCarTypeDetect,
	//processVehicleDamage,
	//processSelfieAnime,
}

var allFuncDescTip = getAllFuncDescWithHead()

var wechatUserToFuncCache = cache.New(24*time.Hour, 1*time.Hour)

var msgIdToResp = cache.New(10*time.Second, 3*time.Minute)

const Token = "wxToken"
const EncodingAESKey = "iRy1YCP2741ZdVSwbnvErub6QZc8uh07SzVfujfIkvq"

type Resp struct {
	mutex   *sync.Mutex
	content *[]byte
}

func getAllFuncDesc() string {
	sb := strings.Builder{}
	for i, s := range allFuncDesc {
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(" : ")
		sb.WriteString(s)
		sb.WriteString("\n")
	}
	return sb.String()
}

func getAllFuncDescWithHead() string {
	sb := strings.Builder{}
	sb.WriteString("功能列表如下、请输入编号切换功能:\n")
	sb.WriteString(getAllFuncDesc())
	sb.WriteString("\n0 : 查看功能列表")
	return sb.String()
}

func processCarTypeDetect(inMessage InMessage) (*OutMessage, error) {
	picUrl := inMessage.PicUrl
	util.Log.Info(picUrl)
	json, err := baidu.CarType("", picUrl, "5", "")
	if err != nil {
		return nil, err
	}
	if json == nil {
		return toTextResponse(inMessage, "调用失败"), nil
	}
	sb := strings.Builder{}
	result := json.Get("result")
	if result.Exists() {
		for _, r := range result.Array() {
			score := fmt.Sprintf("%.2f%%", r.Get("score").Float()*10000000/100000)
			year := r.Get("year").String()
			name := r.Get("name").String()
			sb.WriteString("车型: ")
			sb.WriteString(name)
			sb.WriteString("\n年份: ")
			sb.WriteString(year)
			sb.WriteString("\n可信度: ")
			sb.WriteString(score)
			baike_info := r.Get("baike_info")
			if baike_info.Exists() {
				baike_url := baike_info.Get("baike_url").String()
				image_url := baike_info.Get("image_url").String()
				description := baike_info.Get("description").String()
				sb.WriteString("\n百度百科地址: ")
				sb.WriteString(baike_url)
				sb.WriteString("\n百度百科图片地址: ")
				sb.WriteString(image_url)
				sb.WriteString("\n详情: \n")
				sb.WriteString(description)
			}
			sb.WriteString("\n---------------------------------\n")
		}
	}
	return toTextResponse(inMessage, sb.String()), nil
}

func processVehicleDamage(inMessage InMessage) (*OutMessage, error) {
	picUrl := inMessage.PicUrl
	json, err := baidu.VehicleDamage("", picUrl)
	if err != nil {
		return nil, err
	}
	if json == nil {
		return toTextResponse(inMessage, "调用失败"), nil
	}
	sb := strings.Builder{}
	result := json.Get("result")
	if result.Exists() {
		damageInfo := result.Get("damage_info")
		if damageInfo.Exists() {
			for _, e1 := range damageInfo.Array() {
				parts := e1.Get("parts")
				probability := e1.Get("probability")
				e_type := e1.Get("type")
				sb.WriteString(fmt.Sprintf("位置[%s]、损伤[%s]、概率[%d]", parts.Str, e_type.Str, probability.Int()))
				numericInfo := e1.Get("numeric_info")
				if numericInfo.Exists() {
					sb.WriteString("损伤详细信息:\n")
					for _, e2 := range numericInfo.Array() {
						length := e2.Get("length")
						width := e2.Get("width")
						area := e2.Get("area")
						ratio := e2.Get("ratio")
						sb.WriteString(fmt.Sprintf("[长:(%.2f)cm、宽:(%.2f)cm、面积:(%.2f)cm、损伤相对于部件占比(%.2f%%)]\n", length.Float(), width.Float(), area.Float(), ratio.Float()*100))
					}
				}
			}
		} else {
			sb.WriteString(result.Get("description").String())
		}
	}
	return toTextResponse(inMessage, sb.String()), nil
}

func processSelfieAnime(inMessage InMessage) (*OutMessage, error) {
	picUrl := inMessage.PicUrl
	json, err := baidu.SelfieAnime("", picUrl, "5", "")
	if err != nil {
		return nil, err
	}
	if json == nil {
		return toTextResponse(inMessage, "调用失败"), nil
	}
	image := json.Get("image")
	if image.Exists() {
		decodeString, err := base64.StdEncoding.DecodeString(image.Str)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		mediaId, err := MediaUpload("image", "temp.png", decodeString)
		if err != nil {
			return nil, err
		}
		return toImageResponse(inMessage, mediaId), nil
	} else {
		util.Log.Errorf("receive message%s\n", json.Str)
		return toTextResponse(inMessage, "调用失败"), nil
	}

}

func ProcessInMessage(inMessage InMessage) (*OutMessage, error) {
	fromUserName := inMessage.FromUserName
	val, inFunc := wechatUserToFuncCache.Get(fromUserName)
	var funcIndex int
	if inFunc {
		funcIndex = val.(int)
	}
	switch inMessage.MsgType {
	case "text":
		{
			inputNum, err := strconv.ParseInt(inMessage.Content, 10, 8)
			if err == nil {
				sb := strings.Builder{}
				if inputNum == 0 {
					if inFunc {
						sb.WriteString(fmt.Sprintf("当前处于功能[%s]下\n", allFuncDesc[funcIndex]))
					} else {
						sb.WriteString("当前未选择任何功能\n")
					}
					sb.WriteString(allFuncDescTip)
					return toTextResponse(inMessage, sb.String()), nil
				} else {
					if int(inputNum) > len(allFunc) {
						return toTextResponse(inMessage, fmt.Sprintf("功能编号[%d]无法识别\n", inputNum)+allFuncDescTip), nil
					} else {
						wechatUserToFuncCache.SetDefault(fromUserName, int(inputNum-1))
						return toTextResponse(inMessage, fmt.Sprintf("欢迎进入功能[%s]\n%s", allFuncDesc[inputNum-1], allFuncTip[inputNum-1])), nil
					}
				}
			} else {
				if inFunc {
					return allFunc[funcIndex](inMessage)
				} else {
					return toTextResponse(inMessage, "还未选择功能\n"+allFuncDescTip), nil
				}
			}
		}
	default:
		if inFunc {
			return allFunc[funcIndex](inMessage)
		} else {
			return toTextResponse(inMessage, "还未选择功能\n"+allFuncDescTip), nil
		}
	}

}

func Route(engine *gin.Engine) {
	engine.Any("/wx/handle", process)
}

func process(_ctx *gin.Context) {
	if _ctx.Request.Method == "GET" {
		signature := _ctx.Query("signature")
		timestamp := _ctx.Query("timestamp")
		nonce := _ctx.Query("nonce")
		echostr := _ctx.Query("echostr")
		arr := []string{Token, timestamp, nonce}
		sort.Strings(arr)
		join := strings.Join(arr, "")
		o := sha1.New()
		o.Write([]byte(join))
		sum := hex.EncodeToString(o.Sum([]byte(nil)))
		if signature == sum {
			_, err := _ctx.Writer.WriteString(echostr)
			if err != nil {
				util.Log.Errorf("%+v", err)
			}
		} else {
			util.Log.Warnf("signature not map \nsignature: %s\ntimestamp: %s\nnonce: %s\nechostr: %s\njoin: %s\nsum: %s", signature, timestamp, nonce, echostr, join, sum)
		}
	} else {
		var inMessage InMessage
		err := _ctx.BindXML(&inMessage)
		if err != nil {
			util.Log.Errorf("%+v", err)
		}
		msgId := inMessage.MsgId
		resp, b := msgIdToResp.Get(msgId)
		if b {
			temp := resp.(*Resp)
			func() {
				temp.mutex.Lock()
				defer temp.mutex.Unlock()
				_, err = _ctx.Writer.Write(*temp.content)
				if err != nil {
					util.Log.Errorf("%+v", err)
				}
			}()
		} else {
			temp := Resp{
				mutex:   &sync.Mutex{},
				content: nil,
			}
			func() {
				temp.mutex.Lock()
				defer temp.mutex.Unlock()
				msgIdToResp.SetDefault(msgId, &temp)
				content, err := ProcessInMessage(inMessage)
				if err != nil {
					util.Log.Errorf("%+v", err)
					return
				}
				marshal, _ := xml.Marshal(content)
				temp.content = &marshal
				_, err = _ctx.Writer.Write(marshal)
				if err != nil {
					util.Log.Errorf("%+v", err)
					return
				}
			}()
		}
	}
}

func toTextResponse(inMessage InMessage, content string) *OutMessage {
	outMessage := OutMessage{}
	outMessage.MsgType = "text"
	outMessage.FromUserName = inMessage.ToUserName
	outMessage.ToUserName = inMessage.FromUserName
	outMessage.CreateTime = time.Now().UnixMilli() / 1000
	outMessage.Content = content
	return &outMessage
}

func toImageResponse(inMessage InMessage, mediaId string) *OutMessage {
	outMessage := OutMessage{}
	outMessage.MsgType = "image"
	outMessage.FromUserName = inMessage.ToUserName
	outMessage.ToUserName = inMessage.FromUserName
	outMessage.CreateTime = time.Now().UnixMilli() / 1000
	outMessage.Image.MediaId = mediaId
	return &outMessage
}

type InMessage struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	MsgId        string `xml:"MsgId"`

	//文本
	Content string `xml:"Content"`

	//图片
	PicUrl  string `xml:"PicUrl"`
	MediaId string `xml:"MediaId"`
}

type OutMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`

	//文本
	Content string `xml:"Content"`

	//图片
	Image OutMessage_Image `xml:"Image"`
}

type OutMessage_Image struct {
	XMLName xml.Name `xml:"Image"`
	MediaId string   `xml:"MediaId"`
}
