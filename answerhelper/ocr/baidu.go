package ocr

import (
	"fmt"
	"strings"
	"github.com/ytongshang/rancune/httpclient"
	"github.com/ytongshang/rancune/answerhelper"
	"github.com/ytongshang/rancune/answerhelper/util"
)

const BaiduOcrGeneralUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/general_basic"

const BaiduOcrAccurateUrl = "https://aip.baidubce.com/rest/2.0/ocr/v1/accurate_basic"

type tokenResp struct {
	AccessToken string `json:"access_token"`
}

var cache = util.NewCache()

func init() {
}

func GetToken() (token string, err error) {
	if token, ok := cache.Get("access_token"); ok && token != "" {
		return token, nil
	}
	tokenStruct := &tokenResp{}
	uri := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=" + answerhelper.AppConf.ApiKey + "&client_secret=" + answerhelper.AppConf.ApiSecret
	err = httpclient.Post(uri).ToJSON(tokenStruct)
	if err != nil {
		return "", err
	}
	token = tokenStruct.AccessToken
	if token == "" {
		return "", fmt.Errorf("%s%v", "get baidu ocr token failed:", err)
	}
	cache.Put("access_token", token)
	return token, nil
}

type wordsResp struct {
	WordsResultNum int32 `json:"words_result_num"`
	WordsResult    []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

func Ocr(token, imageBase64 string) (question string, answer []string, err error) {
	request := httpclient.Post(BaiduOcrGeneralUrl)
	request.Param("access_token", token)
	request.Param("image", imageBase64)
	wordsStruct := &wordsResp{}
	err = request.ToJSON(wordsStruct)
	if err != nil {
		return "", nil, err
	}
	question, answer = processText(wordsStruct)
	return question, answer, nil
}

func processText(result *wordsResp) (question string, answer []string) {
	flag := true
	for _, item := range result.WordsResult {
		word := item.Words
		if flag {
			question += word
		} else {
			answer = append(answer, word)
		}
		if strings.HasSuffix(word, "?") || strings.HasSuffix(word, "？") {
			flag = false
		}
	}
	strings.Replace(question, "?", "", -1)
	strings.Replace(question, "？", "", -1)
	return
}
