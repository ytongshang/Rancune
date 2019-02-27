package main

import (
	"fmt"

	"log"

	"encoding/json"

	"github.com/ytongshang/rancune/jpush"
	"github.com/ytongshang/rancune/jpush/config"
)

const (
	MESSAGE      = "message"
	NOTIFICATION = "notification"
)

func main() {
	C, err := config.InitConfig()
	if err != nil {
		log.Fatal(C)
	}
	log.Println(C)

	//Platform
	var pf jpush.Platform
	if InSlice(C.Platforms, jpush.ALL) {
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

	// Option
	var option *jpush.Option
	var apns = C.Opions["apns_production"]
	if value, ok := apns.(bool); ok {
		option = new(jpush.Option)
		option.ApnsProduction = value
	}

	// Message
	var msg *jpush.Message
	if InSlice(C.PushType, MESSAGE) {
		msg = new(jpush.Message)
		msg.Title = C.Msg.Title
		var content string
		b, err := json.Marshal(C.Msg.Extras)
		if err == nil {
			content = string(b)
		}
		msg.Content = content
	}

	// notification
	var notification *jpush.Notification
	if InSlice(C.PushType, NOTIFICATION) {
		notification = new(jpush.Notification)
		notification.SetAlert(C.Notice.Alert)

		// android
		var androidNotification jpush.AndroidNotification
		androidNotification.Alert = C.Notice.Alert
		androidNotification.Title = C.Notice.AndroidTitle
		androidNotification.UriAction = C.Notice.AndroidUriAction
		androidNotification.UriActivity = C.Notice.AndroidUriActivity
		if androidNotification.Extras == nil {
			androidNotification.Extras = make(map[string]interface{})
		}
		b, err := json.Marshal(C.Notice.Extras)
		if err == nil {
			androidNotification.Extras["AndroidPushContent"] = string(b)
		}
		notification.SetAndroidNotice(&androidNotification)

		// ios
		var iosNotification jpush.IOSNotification
		iosNotification.Alert = C.Notice.Alert
		iosNotification.Sound = C.Notice.IOSSound
		iosNotification.Badge = C.Notice.IOSBadge
		iosNotification.Category = C.Notice.IOSCategory
		iosNotification.ContentAvailable = C.Notice.IOSContentAvailable
		iosNotification.MutableContent = C.Notice.IOSMutableContent
		if iosNotification.Extras == nil {
			iosNotification.Extras = make(map[string]interface{})
		}
		b1, err1 := json.Marshal(C.Notice.Extras)
		if err1 == nil {
			iosNotification.Extras["IOSPushContent"] = string(b1)
		}
		notification.SetIOSNotice(&iosNotification)
	}

	if msg == nil && notification == nil {
		fmt.Println("请指定pushtype")
		return
	}

	payload := jpush.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	if option != nil {
		payload.SetOptions(option)
	}
	if msg != nil {
		payload.SetMessage(msg)
	}
	if notification != nil {
		payload.SetNotification(notification)
	}

	bytes, err := payload.ToBytes()
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%s\r\n", string(bytes))

	//push
	c := jpush.NewPushClient(C.AppKey, C.MasterSecret, C.PushBaseUrl)
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
