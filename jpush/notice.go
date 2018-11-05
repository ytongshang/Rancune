package jpush

type Notification struct {
	Alert   string               `json:"alert,omitempty"`
	Android *AndroidNotification `json:"android,omitempty"`
	IOS     *IOSNotification     `json:"ios,omitempty"`
}

type AndroidNotification struct {
	Alert     string                 `json:"alert"`
	Title     string                 `json:"title,omitempty"`
	BuilderId int                    `json:"builder_id,omitempty"`
	UriAction string                 `json:"uri_action,omitempty"`
	Extras    map[string]interface{} `json:"extras,omitempty"`
}

type IOSNotification struct {
	Alert            interface{}            `json:"alert"`
	Sound            string                 `json:"sound,omitempty"`
	Badge            string                 `json:"badge,omitempty"`
	ContentAvailable bool                   `json:"content-available,omitempty"`
	MutableContent   bool                   `json:"mutable-content,omitempty"`
	Category         string                 `json:"category,omitempty"`
	Extras           map[string]interface{} `json:"extras,omitempty"`
}

func (this *Notification) SetAlert(alert string) {
	this.Alert = alert
}

func (this *Notification) SetAndroidNotice(android *AndroidNotification) {
	this.Android = android
}

func (this *Notification) SetIOSNotice(ios *IOSNotification) {
	this.IOS = ios
}
