package baidu

import (
	"bcd_go/util"
	"bytes"
	"github.com/dop251/goja"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"sync"
	"time"
)
import "github.com/tidwall/gjson"

const clientId = "sZnWtPyo8VpnG3TPVy6pIgYg"
const clientSecret = "bI7HokxcdbvMfpY0I3mL6vB2GqsSlxbk"

var globalAccessToken string
var expiredAt int64 = -1
var mutex sync.Mutex

var baiduFanYiSignFn func(msg string) string

func getAccessToken() (string, error) {
	client := resty.New()
	if expiredAt == -1 || time.Now().UnixMilli() >= expiredAt {
		mutex.Lock()
		defer mutex.Unlock()
		if expiredAt == -1 || time.Now().UnixMilli() >= expiredAt {
			post, err := client.R().
				SetQueryParam("grant_type", "client_credentials").
				SetQueryParam("client_id", clientId).
				SetQueryParam("client_secret", clientSecret).
				Post("https://aip.baidubce.com/oauth/2.0/token?grant_type")
			if err != nil {
				return "", errors.WithStack(err)
			}
			body := post.Body()
			code := post.StatusCode()
			//util.Log.Debugf("receive body:\n%s\n", string(body))
			if code == 200 {
				parseBytes := gjson.ParseBytes(body)
				accessToken := parseBytes.Get("access_token")
				if accessToken.Exists() {
					expiresIn := parseBytes.Get("expires_in")
					globalAccessToken = accessToken.Str
					expiredAt = time.Now().UnixMilli() + expiresIn.Int()*1000
					return globalAccessToken, nil
				} else {
					return "", errors.Errorf("get access_token failed,receive error body:\n%s", string(body))
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

func SelfieAnime(image string, url string, p_type string, mask_id string) (*gjson.Result, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"image":   image,
			"url":     url,
			"type":    p_type,
			"mask_id": mask_id,
		}).
		SetQueryParam("access_token", token).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		Post("https://aip.baidubce.com/rest/2.0/image-process/v1/selfie_anime")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("SelfieAnime\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}

}

func CarType(image string, url string, top_num string, baike_num string) (*gjson.Result, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"image":     image,
			"url":       url,
			"top_num":   top_num,
			"baike_num": baike_num,
		}).
		SetQueryParam("access_token", token).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		Post("https://aip.baidubce.com/rest/2.0/image-classify/v1/car")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("CarType\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func VehicleDamage(image string, url string) (*gjson.Result, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"image": image,
			"url":   url,
		}).
		SetQueryParam("access_token", token).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		Post("https://aip.baidubce.com/rest/2.0/image-classify/v1/vehicle_damage")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("VehicleDamage\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func OcrAccurate(image string, url string, pdf_file string, pdf_file_num string, language_type string, detect_direction string, paragraph string, probability string) (*gjson.Result, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"image":            image,
			"url":              url,
			"pdf_file":         pdf_file,
			"pdf_file_num":     pdf_file_num,
			"language_type":    language_type,
			"detect_direction": detect_direction,
			"paragraph":        paragraph,
			"probability":      probability,
		}).
		SetQueryParam("access_token", token).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		Post("https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("OcrAccurate\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func OcrFanyi(from string, imageContent *[]byte) (*gjson.Result, error) {
	client := resty.New()
	post, err := client.R().
		SetCookie(&http.Cookie{Name: "BAIDUID", Value: "1671FA7BA5F61F56711C3C7E7F1776BE:FG=1"}).
		SetFileReader("image", "image", bytes.NewReader(*imageContent)).
		SetFormData(map[string]string{
			"from": from,
			"to":   "zh",
		}).
		Post("https://fanyi.baidu.com/getocr")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("OcrFanyi\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func OcrFormAsync(image string, is_sync string, request_type string) (*gjson.Result, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"image":        image,
			"is_sync":      is_sync,
			"request_type": request_type,
		}).
		SetQueryParam("access_token", token).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		Post("https://aip.baidubce.com/rest/2.0/solution/v1/form_ocr/request")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("OcrFormAsync\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func OcrFormAsyncResult(request_id string, request_type string) (*gjson.Result, error) {
	token, err := getAccessToken()
	if err != nil {
		return nil, err
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"request_id":   request_id,
			"request_type": request_type,
		}).
		SetQueryParam("access_token", token).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		Post("https://aip.baidubce.com/rest/2.0/solution/v1/form_ocr/get_request_result")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("OcrFormAsyncResult\n%s", string(body))
		res := gjson.ParseBytes(body)
		if res.Get("error_code").Int() == 0 {
			return &res, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func Fanyi(query string, from string, to string) (*gjson.Result, error) {
	if baiduFanYiSignFn == nil {
		initBaiduFanYiJs()
	}
	client := resty.New()
	post, err := client.R().
		SetFormData(map[string]string{
			"query":             query,
			"from":              from,
			"to":                to,
			"sign":              baiduFanYiSignFn(query),
			"transtype":         "translang",
			"simple_means_flag": "3",
			"token":             "6e3a577552496c051f11504d13d2c9d5",
			"domain":            "common",
		}).
		SetQueryParam("from", from).
		SetQueryParam("to", to).
		SetHeader("content-type", "application/x-www-form-urlencoded").
		SetHeader("Cookie", "BIDUPSID=36372D03359D23781367EED00C55BF8C; PSTM=1624513480; __yjs_duid=1_9b18d48398b589a299c6339ca21c6c6b1624513484016; BAIDUID=2E1E9EE6E19EB4327FE81BFE224C207B:FG=1; BDUSS=mxseW1qVmJMVFFGZGR3YWJBSjZDVUpyblR1RVlmTFIxSVo0d2VCRjNsOGhXR0ppRVFBQUFBJCQAAAAAAAAAAAEAAACXGX0HYmFpY2hhbmdkYQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACHLOmIhyzpiZ; BDUSS_BFESS=mxseW1qVmJMVFFGZGR3YWJBSjZDVUpyblR1RVlmTFIxSVo0d2VCRjNsOGhXR0ppRVFBQUFBJCQAAAAAAAAAAAEAAACXGX0HYmFpY2hhbmdkYQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACHLOmIhyzpiZ; H_WISE_SIDS=110085_127969_179346_184716_188333_188744_189755_190627_194085_194511_194519_196528_197241_197471_197711_198271_199023_199177_199570_201193_202652_203282_203310_204123_204711_204715_204864_204902_205218_205485_206929_207235_207697_207830_207864_208686_208721_208806_208969_209340_209394_209435_209456_209512_209522_209568_209579_209748_209844_209945_209981_210126_210164_210306_210359_210440_210445_210585_210611_210653_210670_210733_210736_210753_210891_210893_210894_210900_210906_211022_211027_211113_211158_211180_211208_211302_211442_211457_211694_211755_211781_211869_212078_8000075_8000105_8000116_8000128_8000137_8000145_8000150_8000156_8000178_8000179_8000183; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; H_PS_PSSID=36429_36460_36455_31254_36452_35913_36165_36487_36518_36055_36519_36337_26350_36468_36311; delPer=0; PSINO=7; BA_HECTOR=ah050k0ka42l8l81811h9bk6v14; ZFY=2QrutpARs:AcbRGQ8U22GLeS5moHTA4BRlSW1CwKvJQc:C; BCLID=8678925632190422762; BDSFRCVID=_1IOJexroG0xuz7DrowlJ7GLQLweG7bTDYrEkR_0p9ppyA8VJeC6EG0Pts1-dEu-EHtdogKKLgOTHULF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF=tb4qVCL5tI-3jJQvhRO5q4tehHRjbMc9WDTm_D_KKxJJD4QeK65OWxTWQPvI5tovB6cX-pPKKR7OJqLm5xnoyjk7QGQyJRvG3mkjbUOGfn02OIbPef6vWt4syPRrKMRnWNTrKfA-b4ncjRcTehoM3xI8LNj405OTbIFO0KJDJCFKbDDCj5Dbj5PW5ptXt6oL26LX3b7EfKoMOq7_bf--DRDmyN53LU3HtDozaD3TBUjHOUT524bxy5KLDG8q2-RZaK3X0x57yj6UEq6HQT3myP5bbN3i-4jTKmj2Wb3cWhRJ8UbS3fvPBTD02-nBat-OQ6npaJ5nJq5nhMJmb67JD-50eGLsKtoXMCKX3JjV5PK_Hn7zepPayf4pbq7H2M-jMTTE_hb-a-3SOC3jXhQd5PL1jU5n0pcH3C3NMKJMf4J2JqD93x6qLTKkQN3T-ntDHCn4L66NWPbfDn3oyUvJXp0n3tvly5jtMgOBBJ0yQ4b4OR5JjxonDh83bG7MJUutfD7H3KCbfCtMMf5; BCLID_BFESS=8678925632190422762; BDSFRCVID_BFESS=_1IOJexroG0xuz7DrowlJ7GLQLweG7bTDYrEkR_0p9ppyA8VJeC6EG0Pts1-dEu-EHtdogKKLgOTHULF_2uxOjjg8UtVJeC6EG0Ptf8g0M5; H_BDCLCKID_SF_BFESS=tb4qVCL5tI-3jJQvhRO5q4tehHRjbMc9WDTm_D_KKxJJD4QeK65OWxTWQPvI5tovB6cX-pPKKR7OJqLm5xnoyjk7QGQyJRvG3mkjbUOGfn02OIbPef6vWt4syPRrKMRnWNTrKfA-b4ncjRcTehoM3xI8LNj405OTbIFO0KJDJCFKbDDCj5Dbj5PW5ptXt6oL26LX3b7EfKoMOq7_bf--DRDmyN53LU3HtDozaD3TBUjHOUT524bxy5KLDG8q2-RZaK3X0x57yj6UEq6HQT3myP5bbN3i-4jTKmj2Wb3cWhRJ8UbS3fvPBTD02-nBat-OQ6npaJ5nJq5nhMJmb67JD-50eGLsKtoXMCKX3JjV5PK_Hn7zepPayf4pbq7H2M-jMTTE_hb-a-3SOC3jXhQd5PL1jU5n0pcH3C3NMKJMf4J2JqD93x6qLTKkQN3T-ntDHCn4L66NWPbfDn3oyUvJXp0n3tvly5jtMgOBBJ0yQ4b4OR5JjxonDh83bG7MJUutfD7H3KCbfCtMMf5; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1653985508; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1653985872; APPGUIDE_10_0_2=1; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; ab_sr=1.0.1_MDkzYjIzZDhjZWE1NGQxODE3MjVjZDNlZDAxZTQ1ZjJmYjRjMDQzMmY4MjNmZDU2NWE3ZDI4ZTY5ZmY0MGE2NTQyNTRlMDlkNTg1N2Y5OTgzMDMyMjY4NmRhY2E1NmRkNmRmYTM5Zjc3YzgyNDAxZjY2YTc3NDExYmYwYjI5MDAzY2Y4YjgyM2U3MWZhZmY5MGZjMzdkN2FkZjUxZTM2ZTRlMzNhYTIzYWYwNTIwYzM5YjlhY2QyZmY5NjlkMDJk").
		Post("https://fanyi.baidu.com/v2transapi")

	if err != nil {
		return nil, errors.WithStack(err)
	}
	body := post.Body()
	code := post.StatusCode()
	if code == 200 {
		util.Log.Infof("Fanyi\n%s", string(body))
		res := gjson.ParseBytes(body)
		return &res, nil
	} else {
		return nil, errors.Errorf("response code %d", code)
	}
}

func initBaiduFanYiJs() {
	vm := goja.New()
	file, err := os.ReadFile("baidu/baidufanyi.js")
	util.LogErrorStack(err)
	_, err = vm.RunString(string(file))
	util.LogErrorStack(err)
	err = vm.ExportTo(vm.Get("e"), &baiduFanYiSignFn)
	util.LogErrorStack(err)
}
