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
	Title        string                 `json:"title"`
	Content      string                 `json:"content"`
	Extras       map[string]interface{} `json:"extras"`
	Opions       map[string]interface{} `json:"options"`
	Msg          Message                `json:"message"`
	Notice       Notice                 `json:"notification"`
}

type Message struct {
	Itemid       string `json:"itemid"`
	Itemtype     string `json:"itemtype"`
	Itemtitle    string `json:"itemtitle"`
	Messagetitle string `json:"messagetitle"`
	Messagedesc  string `json:"messagedesc"`
	Createdate   string `json:"createdate"`
	Linkurl      string `json:"linkurl"`
	Key          string `json:"key"`
	Ss           string `json:"_ss"`
	FromSource   string `json:"_fromSource"`
}

type Notice struct {
	Alert  string                 `json:"alert"`
	Title  string                 `json:"title"`
	Extras map[string]interface{} `json:"extras"`
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
