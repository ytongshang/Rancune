package answerhelper

import (
	"encoding/json"
	"io/ioutil"

	"log"

	"github.com/fatih/color"
)

const (
	ImagePath = "./answerhelper/images/"
)

var AppConf *AppConfig

type AppConfig struct {
	ApiKey    string   `json:"baiduApiKey"`
	ApiSecret string   `json:"baiduAppSecret"`
	Left      int      `json:"left"`
	Top       int      `json:"top"`
	Right     int      `json:"right"`
	Bottom    int      `json:"bottom"`
	Search    []string `json:"search"`
}

func InitConfig() {
	path := "./answerhelper/config.json"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red("open json failed, %v", err)
		return
	}
	AppConf = &AppConfig{}
	err = json.Unmarshal(content, AppConf)
	if err != nil {
		color.Red("xml json failed,%v", err)
		return
	}
	log.Println(AppConf)
}
