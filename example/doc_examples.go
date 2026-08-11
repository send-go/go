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
