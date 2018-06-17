package main

import (
	"fmt"

	"log"

	"encoding/json"
	"github.com/ytongshang/rancune/jpush"
)

const All = "all"

func main() {
	C, err := initConfig()
	if err != nil {
		log.Fatal(C)
	}
	log.Println(C)

	//Platform
	var pf jpush.Platform
	if InSlice(C.Platforms, All) {
		pf.All()
	} else {
		if InSlice(C.Platforms, jpush.IOS) {
			pf.AddIOS()
		}
		if InSlice(C.Platforms, jpush.ANDROID) {
			pf.AddAndrid()
		}
		if InSlice(C.Platforms, jpush.WINPHONE) {
			pf.AddWinphone()
		}
	}

	//Audience
	var ad jpush.Audience
	if len(C.Tags) > 0 {
		ad.SetTag(C.Tags)
	}
	if len(C.TagsNot) > 0 {
		ad.SetTagNot(C.TagsNot)
	}
	if len(C.TagsAnd) > 0 {
		ad.SetTagNot(C.TagsAnd)
	}
	if len(C.Ids) > 0 {
		ad.SetID(C.Ids)
	}
	if len(C.Alias) > 0 {
		ad.SetAlias(C.Alias)
	}

	// Message
	var msg jpush.Message
	msg.Title = C.Title
	content := C.Content
	if content == "" {
		b, err := json.Marshal(C.Msg)
		if err == nil {
			content = string(b)
		}
	}
	msg.Content = content
	for key, value := range C.Extras {
		msg.AddExtras(key, value)
	}

	payload := jpush.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.SetMessage(&msg)

	bytes, err := payload.ToBytes()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%s\r\n", string(bytes))

	//push
	c := jpush.NewPushClient(C.MasterSecret, C.AppKey, C.PushBaseUrl)
	str, err := c.SendPush(bytes)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	} else {
		fmt.Printf("ok:%s", str)
	}
}

func InSlice(slice []string, value string) bool {
	for _, vv := range slice {
		if vv == value {
			return true
		}
	}
	return false
}
