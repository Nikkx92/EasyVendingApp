package domain

type MetaDetails struct {
	UserAgent string `json:"userAgent"`
}

type DeviceInfo struct {
	SourceDeviceID string      `json:"sourceDeviceId"`
	SourceType     string      `json:"sourceType"`
	AppVersion     string      `json:"appVersion"`
	MetaDetails    MetaDetails `json:"metaDetails"`
}

type AuthData struct {
	CompanyId   string `json:"CompanyId"`
	UserLogin   string `json:"UserLogin"`
	PasswordKit string `json:"PasswordKit"`
	INN         string `json:"INN"`
	PasswordFns string `json:"PasswordFns"`
	TimeOffset  int64  `json:"timeOffset"`
}

type SingleRequest struct {
	Jwt      string   `json:"jwt"`
	Date     string   `json:"date"`
	DataKit  []string `json:"dataKit"`
	DeviceID string   `json:"deviceId"`
	AutoMode bool     `json:"autoMode"`
	IsPaid   bool     `json:"isPaid"`
}

type Result struct {
	Message string
	Err     error
}

type Time struct {
	Zone   string
	Offset int32
}
type AutoModeData struct {
	CompanyID    string
	UserLogin    string
	PasswordKit  string
	INN          string
	PasswordFns  string
	DeviceInfo   *DeviceInfo
	RefreshToken string
	Token        string
	Time         *Time
}

type Sale struct {
	DateTime  string  `json:"DateTime"`
	GoodsName string  `json:"GoodsName"`
	Sum       float64 `json:"Sum"`
}
