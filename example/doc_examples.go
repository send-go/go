// Compiles the API surface used by the published guides.
// Never run — `go build` proving it type-checks is the point.
package main

import (
	"errors"
	"log"
	"os"

	"github.com/send-go/go/sendgo"
)

func alertOps(msg string) { log.Println(msg) }

func main() {
	client, err := sendgo.New(sendgo.Config{
		AccessKey:      os.Getenv("SENDGO_ACCESS_KEY"),
		SecretKey:      os.Getenv("SENDGO_SECRET_KEY"),
		KakaoSenderKey: os.Getenv("SENDGO_KAKAO_SENDER_KEY"),
		SmsSenderKey:   os.Getenv("SENDGO_SMS_SENDER_KEY"),
		APIVersion:     "v2",
	})
	if err != nil {
		log.Fatal(err)
	}

	// --- alimtalk (pointer fields via sendgo.String) ---
	_ = client.Alimtalk.Send(sendgo.AlimtalkRequest{
		TemplateCode: "ORDER_CONFIRM_001",
		Contacts: []sendgo.Contact{
			{Contact: "01011111111", Name: "Hong", Var1: "ORD-001", Var8: "x"},
		},
	})
	_ = client.Alimtalk.Send(sendgo.AlimtalkRequest{
		TemplateCode: "PROMO",
		ScheduleType: "SCHEDULED",
		At:           sendgo.String("2026-07-28 09:00:00"),
		Contacts:     []sendgo.Contact{{Contact: "01012345678"}},
	})
	_ = client.Alimtalk.Send(sendgo.AlimtalkRequest{
		TemplateCode: "DELIVERY",
		ReplaceSms:   "Y",
		SmsSubject:   sendgo.String("[shipping]"),
		SmsContent:   sendgo.String("shipped"),
		Contacts:     []sendgo.Contact{{Contact: "01012345678"}},
	})

	// --- friendtalk ---
	_ = client.Friendtalk.Send(sendgo.FriendtalkRequest{
		Content:  "hello",
		Contacts: []sendgo.Contact{{Contact: "01012345678"}},
	})

	// --- brand message ---
	result, err := client.BrandMessage.Send(sendgo.BrandMessageRequest{
		Targeting:          "M",
		MessageType:        "FL",
		FriendTemplateUUID: "9cd5460b-6458-4edc-9b11-c26d3013c340",
		Contacts:           []sendgo.Contact{{Contact: "01012345678", Var1: "29,000"}},
	})
	if err != nil {
		log.Fatal(err)
	}
	_, _ = client.BrandMessage.Broadcast(sendgo.BrandMessageRequest{
		MessageType:        "FW",
		FriendTemplateUUID: "9cd5460b-6458-4edc-9b11-c26d3013c340",
	})
	campaignID, _ := result["data"].(map[string]any)["campaignId"].(string)
	_, _ = client.BrandMessage.Campaign(campaignID)
	_, _ = client.BrandMessage.Campaigns(sendgo.BrandMessageListQuery{From: "2026-08-01", Count: 10})

	// --- sms ---
	_ = client.SMS.SendSMS(sendgo.SmsRequest{
		Content:  "code 123456",
		Contacts: []sendgo.Contact{{Contact: "01012345678"}},
	})
	_ = client.SMS.SendLMS(sendgo.SmsRequest{
		Subject:  sendgo.String("[notice]"),
		Content:  "long text",
		Contacts: []sendgo.Contact{{Contact: "01012345678"}},
	})
	_ = client.SMS.SendMMS(sendgo.SmsRequest{
		Subject:  sendgo.String("[event]"),
		Content:  "deals",
		Contacts: []sendgo.Contact{{Contact: "01012345678"}},
	})

	// --- short URL ---
	created, err := client.ShortURL.Create(sendgo.ShortURLRequest{
		TargetURL: "https://example.com/promotions/summer-sale",
		Title:     "Summer sale landing",
	})
	if err != nil {
		log.Fatal(err)
	}
	data := created["data"].(map[string]any)
	code := data["code"].(string)

	_, _ = client.ShortURL.Stats(code, sendgo.ShortURLStatsQuery{From: "2026-08-01"})
	_, _ = client.ShortURL.List(sendgo.ShortURLListQuery{Count: 10})
	_, _ = client.ShortURL.Show(code)
	_, _ = client.ShortURL.Deactivate(code)

	// --- error handling ---
	var sgErr *sendgo.SendgoError
	if errors.As(err, &sgErr) {
		log.Printf("Sendgo %d [%s]: %s", sgErr.StatusCode, sgErr.ErrorCode, sgErr.Message)
		switch sgErr.ErrorCode {
		case "INVALID_ACCESS_KEY", "INVALID_SECRET_KEY":
			alertOps("check keys")
		case "IP_NOT_ALLOWED":
			alertOps("ip not allowed")
		case "PAYMENT_REQUIRED":
			alertOps("out of credit")
		}
	}
}

// managementApi compiles the 1.3.0 관리 API surface used by the guides.
//
// 등록·심사는 발송과 달리 즉시 완료되지 않는다. 인증번호는 채널 관리자
// 휴대폰으로 가므로 여기서는 발송 트리거까지만 쓴다 — 코드 제출은 여러분
// 화면에서 받아 Create 로 넘긴다.
func managementApi(client *sendgo.Client) {
	// --- 카카오 채널 등록 (2단계) ---
	_, _ = client.KakaoSenders.RequestToken("@my-channel", "01012345678")

	created, err := client.KakaoSenders.Create(sendgo.KakaoSenderCreateRequest{
		Token:        "123456",
		YellowID:     "@my-channel",
		PhoneNumber:  "01012345678",
		CategoryCode: "001001",
	})
	if err != nil {
		log.Fatal(err)
	}

	sender := created["data"].(map[string]any)["sender"].(map[string]any)
	kakaoSenderKey := sender["kakaoSenderKey"].(string)

	_, _ = client.KakaoSenders.Categories("")
	_, _ = client.KakaoSenders.List()
	_, _ = client.KakaoSenders.Sync("")
	_, _ = client.KakaoSenders.Sync(kakaoSenderKey)
	_, _ = client.KakaoSenders.ApplyBrandMessageTargeting(kakaoSenderKey, "N")

	// --- 알림톡 템플릿 등록 → 검수 요청 → 폴링 ---
	template, err := client.NoticeTemplates.Create(sendgo.NoticeTemplateRequest{
		KakaoSenderKey:        kakaoSenderKey,
		TemplateName:          "주문 접수 안내",
		TemplateContent:       "#{name}님, 주문 #{orderNo}이 접수되었습니다.",
		TemplateMessageType:   "BA",
		TemplateEmphasizeType: "NONE",
		CategoryCode:          "001001",
		MessagePurpose:        "order_delivery",
		LegalBasis:            "transaction",
		BenefitOrigin:         "none",
		ExpiryType:            "none",
		OptInReviewConfirmed:  true,
		CtaClearConfirmed:     true,
		PolicyConfirmed:       true,
	})
	if err != nil {
		log.Fatal(err)
	}

	templateCode := template["data"].(map[string]any)["template"].(map[string]any)["templateCode"].(string)

	_, _ = client.NoticeTemplates.RequestInspection(templateCode, "", nil)
	_, _ = client.NoticeTemplates.Sync(templateCode)
	_, _ = client.NoticeTemplates.List(sendgo.NoticeTemplateListQuery{
		KakaoSenderKey:   kakaoSenderKey,
		InspectionStatus: "APR",
	})
	_, _ = client.NoticeTemplates.CancelInspection(templateCode)
	_, _ = client.NoticeTemplates.Release(templateCode)
	_, _ = client.NoticeTemplates.Delete(templateCode)

	// --- 이미지 템플릿 / 검수 첨부 (multipart) ---
	image, err := os.Open("banner.jpg")
	if err == nil {
		defer image.Close()
		_, _ = client.NoticeTemplates.CreateWithImage(
			sendgo.NoticeTemplateRequest{
				KakaoSenderKey:        kakaoSenderKey,
				TemplateName:          "이벤트 안내",
				TemplateContent:       "#{name}님께 드리는 안내입니다.",
				TemplateEmphasizeType: "IMAGE",
				CategoryCode:          "001001",
				MessagePurpose:        "service_ops",
				LegalBasis:            "transaction",
				BenefitOrigin:         "none",
				ExpiryType:            "none",
				OptInReviewConfirmed:  true,
				CtaClearConfirmed:     true,
				PolicyConfirmed:       true,
			},
			sendgo.MultipartFile{FileName: "banner.jpg", Content: image},
		)
	}

	// --- 브랜드메시지 템플릿 ---
	_, _ = client.BrandTemplates.Create(sendgo.BrandTemplateRequest{
		KakaoSenderKey:  kakaoSenderKey,
		TemplateName:    "여름 세일 안내",
		TemplateType:    "FI",
		TemplateContent: "여름 세일이 시작되었습니다.",
		ImageURL:        "https://mud-kage.kakao.com/example.jpg",
	})
	_, _ = client.BrandTemplates.List(sendgo.BrandTemplateListQuery{KakaoSenderKey: kakaoSenderKey})
	_, _ = client.BrandTemplates.Import(kakaoSenderKey)

	// --- 발신번호 등록 신청 ---
	_, _ = client.SenderRegistration.NumberTypes()
	_, _ = client.SenderRegistration.Validate("02-1234-5678", "team_main")

	csu, err := os.Open("csu.pdf")
	if err == nil {
		defer csu.Close()
		_, _ = client.SenderRegistration.Create(
			sendgo.SenderRegistrationRequest{
				SenderAlias:      "고객센터 대표번호",
				SenderNumberType: "team_main",
				PhoneE164:        "02-1234-5678",
			},
			[]sendgo.MultipartFile{
				{FieldName: "csuCertificate", FileName: "csu.pdf", Content: csu},
			},
		)
	}

	_, _ = client.SenderRegistration.List()
	_, _ = client.SenderRegistration.Update("sender-key", "새 이름", "")

	// --- 문자 상용구 템플릿 ---
	_, _ = client.MessageTemplates.Create(sendgo.MessageTemplateRequest{
		MessageTranType:    "LMS",
		MessageTranSubject: "주문 안내",
		MessageTranMsg:     "주문이 접수되었습니다.",
	})
	_, _ = client.MessageTemplates.List(sendgo.MessageTemplateListQuery{MessageType: "LMS"})
}
