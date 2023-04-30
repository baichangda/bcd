package wechat

import (
	"bytes"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"sync"
	"time"
)

var globalAccessToken string
var expiredAt int64 = -1
var mutex sync.Mutex

const appid = "wxe3962f308f7be716"
const secret = "31d221988cdc54c6fd9f88bdef3e8e34"

func GetAccessToken() (string, error) {
	if expiredAt == -1 || time.Now().UnixMilli() >= expiredAt {
		mutex.Lock()
		defer mutex.Unlock()
		if expiredAt == -1 || time.Now().UnixMilli() >= expiredAt {
			client := resty.New()
			get, err := client.R().
				SetQueryParam("grant_type", "client_credential").
				SetQueryParam("appid", appid).
				SetQueryParam("secret", secret).
				Get("https://api.weixin.qq.com/cgi-bin/token")
			if err != nil {
				return "", errors.WithStack(err)
			}
			body := get.Body()
			code := get.StatusCode()
			//cfg.Log.Debugf("receive body:\n%s\n", string(body))
			if code == 200 {
				parseBytes := gjson.ParseBytes(body)
				accessToken := parseBytes.Get("access_token")
				if accessToken.Exists() {
					expiresIn := parseBytes.Get("expires_in")
					globalAccessToken = accessToken.Str
					expiredAt = time.Now().UnixMilli() + expiresIn.Int()*1000
					return globalAccessToken, nil
				} else {
					return "", errors.Errorf("getAccessToken failed,receive error body:\n%s", string(body))
				}
			} else {
				return "", errors.Errorf("response code %d", code)
			}
		} else {
			return globalAccessToken, nil
		}
	} else {
		return globalAccessToken, nil
	}

}

func MediaUpload(p_type string, fileName string, fileContent []byte) (string, error) {
	token, err := GetAccessToken()
	if err != nil {
		return "", err
	}
	client := resty.New()
	post, err := client.R().
		SetQueryParam("access_token", token).
		SetQueryParam("type", p_type).
		SetFileReader("media", fileName, bytes.NewReader(fileContent)).
		Post("https://api.weixin.qq.com/cgi-bin/media/upload")
	if err != nil {
		return "", errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	//util.Log.Debugf("receive body:\n%s\n", string(body))
	if code == 200 {
		parseBytes := gjson.ParseBytes(body)
		mediaId := parseBytes.Get("media_id")
		if mediaId.Exists() {
			return mediaId.Str, nil
		} else {
			return "", errors.Errorf("mediaUpload failed,receive error body:\n%s", string(body))
		}
	} else {
		return "", errors.Errorf("response code %d", code)
	}
}
