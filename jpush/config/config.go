package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	AppKey       string                 `json:"appkey"`
	MasterSecret string                 `json:"mastersecret"`
	PushBaseUrl  string                 `json:"pushbaseurl"`
	Platforms    []string               `json:"platforms"`
	Alias        []string               `json:"alias"`
	Tags         []string               `json:"tags"`
	TagsAnd      []string               `json:"tag_and"`
	TagsNot      []string               `json:"tag_not"`
	Ids          []string               `json:"ids"`
	Opions       map[string]interface{} `json:"options"`
	PushType     []string               `json:"pushtype"`
	Msg          Message                `json:"message"`
	Notice       Notice                 `json:"notification"`
}

type Message struct {
	Title  string                 `json:"title"`
	Extras map[string]interface{} `json:"extras"`
}

type Notice struct {
	Alert               string                 `json:"alert"`
	AndroidTitle        string                 `json:"android-title"`
	AndroidUriAction    string                 `json:"android-uri-action"`
	IOSSound            string                 `json:"ios-sound"`
	IOSBadge            string                 `json:"ios-badge"`
	IOSContentAvailable bool                   `json:"ios-content-available"`
	IOSMutableContent   bool                   `json:"ios-mutable-content"`
	IOSCategory         string                 `json:"ios-category"`
	Extras              map[string]interface{} `json:"extras"`
}

func InitConfig() (*config, error) {
	path := "./pushconfig/jpushconfig.json"
	filename, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	conf := &config{}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
