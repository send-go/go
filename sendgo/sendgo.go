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
	Alimtalk *AlimtalkService
	// Deprecated: 친구톡은 2025-12-31 종료. BrandMessage를 사용하세요.
	Friendtalk *FriendtalkService
	// BrandMessage는 카카오 브랜드메시지(친구톡의 후속 채널)입니다. v2 전용.
	BrandMessage *BrandMessageService
	// ShortURL은 짧은 URL — 링크 단축과 클릭 반응 분석입니다. v2 전용.
	ShortURL *ShortURLService
	SMS      *SMSService

	// ---------------------------------------------------- 관리 API (v2 전용)
	// 콘솔에서만 되던 등록·심사. 발송과 달리 대부분 즉시 완료되지 않습니다 —
	// 등록 성공은 "접수됨"이지 "사용 가능"이 아닙니다.

	// KakaoSenders는 카카오 발신프로필(채널) 등록·동기화입니다. 기업 계정 전용.
	KakaoSenders *KakaoSenderService
	// NoticeTemplates는 알림톡 템플릿 등록·수정·검수 요청입니다. 기업 계정 전용.
	NoticeTemplates *NoticeTemplateService
	// BrandTemplates는 브랜드메시지(구 친구톡) 템플릿 관리입니다. 기업 계정 전용.
	BrandTemplates *BrandTemplateService
	// SenderRegistration은 발신번호 등록·심사 접수입니다.
	SenderRegistration *SenderRegistrationService
	// MessageTemplates는 문자 상용구 템플릿입니다.
	MessageTemplates *MessageTemplateService
	// KakaoImages는 카카오 이미지 업로드입니다. 기업 계정 전용.
	KakaoImages *KakaoImageService
	// RejectedNumbers는 수신거부(080) 번호 조회입니다.
	RejectedNumbers *RejectedNumberService
	// Webhook은 이벤트 웹훅 구독입니다 — 등록·심사 결과를 밀어 받습니다.
	Webhook *WebhookService
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
		Alimtalk:     newAlimtalkService(hc, cfg.KakaoSenderKey, cfg.SmsSenderKey),
		Friendtalk:   newFriendtalkService(hc, cfg.KakaoSenderKey, cfg.SmsSenderKey),
		BrandMessage: newBrandMessageService(hc, cfg.KakaoSenderKey, cfg.SmsSenderKey),
		ShortURL:     newShortURLService(hc),
		SMS:          newSMSService(hc, cfg.SmsSenderKey),

		KakaoSenders:       newKakaoSenderService(hc),
		NoticeTemplates:    newNoticeTemplateService(hc),
		BrandTemplates:     newBrandTemplateService(hc),
		SenderRegistration: newSenderRegistrationService(hc),
		MessageTemplates:   newMessageTemplateService(hc),
		KakaoImages:        newKakaoImageService(hc),
		RejectedNumbers:    newRejectedNumberService(hc),
		Webhook:            newWebhookService(hc),
	}, nil
}
