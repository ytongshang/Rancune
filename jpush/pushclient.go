package jpush

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	HOST_NAME_SSL     = "https://api.jpush.cn/v3/push"
	HOST_SCHEDULE     = "https://api.jpush.cn/v3/schedules"
	HOST_REPORT       = "https://report.jpush.cn/v3/received"
	CHARSET           = "UTF-8"
	CONTENT_TYPE_JSON = "application/json"
)

type PushClient struct {
	MasterSecret string
	AppKey       string
	AuthCode     string
	BaseUrl      string
}

func NewPushClient(appKey, secret, baseurl string) *PushClient {
	base64Str := base64.StdEncoding.EncodeToString([]byte(appKey + ":" + secret))
	auth := "Basic " + base64Str
	if baseurl == "" {
		baseurl = HOST_NAME_SSL
	}
	pusher := &PushClient{secret, appKey, auth, baseurl}
	return pusher
}

func (this *PushClient) SendPush(content []byte) (string, error) {
	return this.SendPushBytes(this.BaseUrl, content)
}

func (this *PushClient) SendPushBytes(url string, content []byte) (string, error) {
	ret, err := sendPostBytes(url, content, this.AuthCode)
	if err != nil {
		return ret, err
	}
	if strings.Contains(ret, "msg_id") {
		return ret, nil
	} else {
		return "", errors.New(ret)
	}
}

func (this *PushClient) CreateSchedule(data []byte) (string, error) {
	return this.SendScheduleBytes(data, HOST_SCHEDULE)
}
func (this *PushClient) SendScheduleBytes(content []byte, url string) (string, error) {
	ret, err := sendPostBytes(url, content, this.AuthCode)
	if err != nil {
		return ret, err
	}
	if strings.Contains(ret, "schedule_id") {
		return ret, nil
	} else {
		return "", errors.New(ret)
	}
}

func (this *PushClient) DeleteSchedule(id string) (string, error) {
	return this.SendDeleteScheduleRequest(id, HOST_SCHEDULE)
}

func (this *PushClient) SendDeleteScheduleRequest(schedule_id string, url string) (string, error) {
	rsp, err := Delete(strings.Join([]string{url, schedule_id}, "/")).Header("Authorization", this.AuthCode).String()
	if err != nil {
		return "", err
	}
	_, err = UnmarshalResponse(rsp)
	if err != nil {
		return "", err
	}
	return rsp, nil
}

func (this *PushClient) GetSchedule(id string) (string, error) {
	return this.SendGetScheduleRequest(id, HOST_SCHEDULE)
}

func (this *PushClient) SendGetScheduleRequest(schedule_id string, url string) (string, error) {
	rsp, err := Get(strings.Join([]string{url, schedule_id}, "/")).Header("Authorization", this.AuthCode).String()
	if err != nil {
		return "", err
	}
	_, err = UnmarshalResponse(rsp)
	if err != nil {
		return "", err
	}
	return rsp, nil
}

func (this *PushClient) GetReport(msg_ids string) (string, error) {
	return this.SendGetReportRequest(msg_ids, HOST_REPORT)
}

func (this *PushClient) SendGetReportRequest(msg_ids string, url string) (string, error) {
	return Get(url).SetBasicAuth(this.AppKey, this.MasterSecret).Param("msg_ids", msg_ids).String()
}

func UnmarshalResponse(rsp string) (map[string]interface{}, error) {
	mapRs := map[string]interface{}{}
	if len(strings.TrimSpace(rsp)) == 0 {
		return mapRs, nil
	}
	err := json.Unmarshal([]byte(rsp), &mapRs)
	if err != nil {
		return nil, err
	}
	if _, ok := mapRs["error"]; ok {
		return nil, fmt.Errorf(rsp)
	}
	return mapRs, nil
}

func sendPostBytes(url string, data []byte, authCode string) (string, error) {
	return Post(url).AddHeader("Charset", CHARSET).
		AddHeader("Authorization", authCode).
		AddHeader("Content-Type", CONTENT_TYPE_JSON).
		Body(data).String()
}
