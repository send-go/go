package sendgo

// Contact는 수신자 정보입니다.
type Contact struct {
	Contact string `json:"contact"`
	Name    string `json:"name,omitempty"`
	Var1    string `json:"var1,omitempty"`
	Var2    string `json:"var2,omitempty"`
	Var3    string `json:"var3,omitempty"`
	Var4    string `json:"var4,omitempty"`
	Var5    string `json:"var5,omitempty"`
	Var6    string `json:"var6,omitempty"`
	Var7    string `json:"var7,omitempty"`
	Var8    string `json:"var8,omitempty"`
}

// AlimtalkRequest는 알림톡 전송 요청입니다.
type AlimtalkRequest struct {
	At           *string   `json:"at"`
	ScheduleType string    `json:"scheduleType"`
	TemplateCode string    `json:"templateCode"`
	ReplaceSms   string    `json:"replaceSms"`
	SmsSubject   *string   `json:"smsSubject"`
	SmsContent   *string   `json:"smsContent"`
	Contacts     []Contact `json:"contacts"`
	// 자동 설정됨
	KakaoSenderKey string `json:"kakaoSenderKey,omitempty"`
	SenderKey      string `json:"senderKey,omitempty"`
}

// FriendtalkRequest는 친구톡 전송 요청입니다.
type FriendtalkRequest struct {
	At           *string   `json:"at"`
	ScheduleType string    `json:"scheduleType"`
	MessageType  string    `json:"messageType"`
	Content      string    `json:"content"`
	Buttons      []any     `json:"buttons"`
	Image        any       `json:"image"`
	ImageURL     *string   `json:"imageUrl"`
	ImageLink    *string   `json:"imageLink"`
	AdFlag       string    `json:"adFlag"`
	Wide         string    `json:"wide"`
	Adult        string    `json:"adult"`
	Header       *string   `json:"header"`
	ReplaceSms   string    `json:"replaceSms"`
	SmsSubject   *string   `json:"smsSubject"`
	SmsContent   *string   `json:"smsContent"`
	Contacts     []Contact `json:"contacts"`
	// 자동 설정됨
	KakaoSenderKey string `json:"kakaoSenderKey,omitempty"`
	SenderKey      string `json:"senderKey,omitempty"`
}

// SmsRequest는 SMS/LMS/MMS 전송 요청입니다.
type SmsRequest struct {
	CampaignType string    `json:"campaignType"`
	MessageType  string    `json:"messageType"`
	ScheduleType string    `json:"scheduleType"`
	At           *string   `json:"at"`
	Subject      *string   `json:"subject"`
	Content      string    `json:"content"`
	Files        []any     `json:"files"`
	Contacts     []Contact `json:"contacts"`
	// 자동 설정됨
	SenderKey string `json:"senderKey,omitempty"`
}

// Config는 Sendgo 클라이언트 설정입니다.
type Config struct {
	AccessKey      string
	SecretKey      string
	KakaoSenderKey string
	SmsSenderKey   string
	APIVersion     string // "v1" | "v2" (기본값: "v1")
	BaseURL        string // 기본값: "https://api.sendgo.io"
}
