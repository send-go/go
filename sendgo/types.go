package sendgo

import "encoding/json"

// Contact는 수신자 정보입니다.
type Contact struct {
	Contact string
	Name    string
	Var1    string
	Var2    string
	Var3    string
	Var4    string
	Var5    string
	Var6    string
	Var7    string
	Var8    string
	// Variables는 임의 명명 템플릿 변수입니다 (예: {"title": "..."} → 알림톡 #{title} 치환).
	// contact 오브젝트에 평탄하게(flat) 직렬화됩니다.
	Variables map[string]string
}

// MarshalJSON은 기본 필드와 Variables 명명 변수를 하나의 contact 오브젝트로 평탄화합니다.
func (c Contact) MarshalJSON() ([]byte, error) {
	m := map[string]string{"contact": c.Contact}
	for k, v := range map[string]string{
		"name": c.Name, "var1": c.Var1, "var2": c.Var2, "var3": c.Var3, "var4": c.Var4,
		"var5": c.Var5, "var6": c.Var6, "var7": c.Var7, "var8": c.Var8,
	} {
		if v != "" {
			m[k] = v
		}
	}
	for k, v := range c.Variables {
		m[k] = v
	}
	return json.Marshal(m)
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
	BaseURL        string // 기본값: "https://sendgo.io"
}
