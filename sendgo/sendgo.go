// Package sendgo는 Sendgo API를 위한 Go SDK입니다.
// 카카오 알림톡/친구톡, SMS/LMS/MMS 전송을 지원합니다.
//
// 사용법:
//
//	client, err := sendgo.New(sendgo.Config{
//	    AccessKey:      os.Getenv("SENDGO_ACCESS_KEY"),
//	    SecretKey:      os.Getenv("SENDGO_SECRET_KEY"),
//	    KakaoSenderKey: os.Getenv("SENDGO_KAKAO_KEY"),
//	    SmsSenderKey:   os.Getenv("SENDGO_SMS_KEY"),
//	    APIVersion:     "v2",
//	})
//
//	err = client.Alimtalk.Send(sendgo.AlimtalkRequest{
//	    TemplateCode: "ORDER_CONFIRM_001",
//	    Contacts:     []sendgo.Contact{{Contact: "01012345678", Var1: "ORD-001"}},
//	})
package sendgo

import "fmt"

// Client는 Sendgo API 클라이언트입니다.
type Client struct {
	Alimtalk   *AlimtalkService
	Friendtalk *FriendtalkService
	SMS        *SMSService
}

// New는 새 Sendgo 클라이언트를 생성합니다.
func New(cfg Config) (*Client, error) {
	if cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("sendgo: AccessKey와 SecretKey는 필수입니다")
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v1"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://sendgo.io"
	}

	tm := newTokenManager(cfg.BaseURL, cfg.AccessKey, cfg.SecretKey, cfg.APIVersion)
	hc := newHTTPClient(tm, cfg.BaseURL, cfg.APIVersion)

	return &Client{
		Alimtalk:   newAlimtalkService(hc, cfg.KakaoSenderKey, cfg.SmsSenderKey),
		Friendtalk: newFriendtalkService(hc, cfg.KakaoSenderKey, cfg.SmsSenderKey),
		SMS:        newSMSService(hc, cfg.SmsSenderKey),
	}, nil
}
