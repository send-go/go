package sendgo

import (
	"encoding/json"
	"strconv"
)

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

// BrandMessageRequest는 카카오 브랜드메시지 전송 요청입니다.
//
// 브랜드메시지는 친구톡의 후속 채널로, MessageType은 친구톡 코드
// (FT/FI/FW/FL/FC/FM/FP/FA)를 그대로 넘기며 브랜드메시지 코드
// (BT/BI/BW/BL/BC/BM/BP/BA) 변환은 서버가 처리합니다.
//
// Targeting은 M(채널 친구) / N(비친구) / I(전체) / F(동보)이며,
// F는 수신자 목록을 카카오 측에서 확장하므로 Contacts를 넘기지 않습니다.
type BrandMessageRequest struct {
	At                 *string  `json:"at"`
	ScheduleType       string   `json:"scheduleType"`
	Targeting          string   `json:"targeting"`
	MessageType        string   `json:"messageType"`
	FriendTemplateUUID string   `json:"friendTemplateUuid"`
	Content            *string  `json:"content"`
	Buttons            []any    `json:"buttons"`
	ImageURL           *string  `json:"imageUrl"`
	ImageLink          *string  `json:"imageLink"`
	AdFlag             string   `json:"adFlag"`
	Adult              string   `json:"adult"`
	PushAlarm          string   `json:"pushAlarm"`
	Header             *string  `json:"header"`
	Coupon             any      `json:"coupon"`
	Item               any      `json:"item"`
	Commerce           any      `json:"commerce"`
	List               []any    `json:"list"`
	Head               any      `json:"head"`
	Tail               any      `json:"tail"`
	Video              any      `json:"video"`
	AdditionalContent  *string  `json:"additionalContent"`
	FriendGroupKey     *string  `json:"friendGroupKey"`
	ReplaceSms         string   `json:"replaceSms"`
	SmsSubject         *string  `json:"smsSubject"`
	SmsContent         *string  `json:"smsContent"`
	RejectServiceID    *string  `json:"rejectServiceId"`
	Webhooks           []string `json:"webhooks"`
	// Contacts는 Targeting이 F(동보)일 때 생략됩니다.
	Contacts []Contact `json:"contacts,omitempty"`
	// 자동 설정됨
	KakaoSenderKey string `json:"kakaoSenderKey,omitempty"`
	SenderKey      string `json:"senderKey,omitempty"`
}

// BrandMessageListQuery는 브랜드메시지 캠페인 목록 조회 조건입니다.
// 빈 문자열/0은 전송되지 않고 서버 기본값이 적용됩니다.
type BrandMessageListQuery struct {
	From  string
	To    string
	Count int
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

// ShortURLRequest는 짧은 URL 생성 요청입니다.
type ShortURLRequest struct {
	// TargetURL은 줄일 원본 URL입니다. http/https 만 허용됩니다.
	TargetURL string `json:"targetUrl"`
	// Title은 관리 화면에서 구분하기 위한 이름입니다.
	Title string `json:"title,omitempty"`
	// ExpiresAt 이후에는 리다이렉트하지 않고 410 Gone 을 반환합니다.
	ExpiresAt string `json:"expiresAt,omitempty"`
	// ForceNew가 true면 같은 URL이라도 새 코드를 만듭니다.
	// 캠페인별로 반응을 분리해 집계할 때 사용합니다.
	ForceNew bool `json:"forceNew,omitempty"`
}

// ShortURLListQuery는 짧은 URL 목록 조회 조건입니다.
type ShortURLListQuery struct {
	From  string
	To    string
	Count int
}

func (q ShortURLListQuery) toMap() map[string]string {
	m := map[string]string{}
	if q.From != "" {
		m["from"] = q.From
	}
	if q.To != "" {
		m["to"] = q.To
	}
	if q.Count > 0 {
		m["count"] = strconv.Itoa(q.Count)
	}
	return m
}

// ShortURLStatsQuery는 짧은 URL 통계 조회 조건입니다.
type ShortURLStatsQuery struct {
	From string
	To   string
}

func (q ShortURLStatsQuery) toMap() map[string]string {
	m := map[string]string{}
	if q.From != "" {
		m["from"] = q.From
	}
	if q.To != "" {
		m["to"] = q.To
	}
	return m
}
