package sendgo

import "strconv"

// BrandMessageService는 카카오 브랜드메시지 전송 서비스입니다.
//
// 브랜드메시지는 친구톡의 후속 채널로, 친구톡과 달리 채널 친구가 아닌
// 수신자에게도 보낼 수 있고(Targeting "N"), 수신 동의한 전체 채널 친구에게
// 동보 발송할 수 있습니다(Targeting "F").
type BrandMessageService struct {
	http           *httpClient
	kakaoSenderKey string
	smsSenderKey   string
}

func newBrandMessageService(hc *httpClient, kakaoKey, smsKey string) *BrandMessageService {
	return &BrandMessageService{http: hc, kakaoSenderKey: kakaoKey, smsSenderKey: smsKey}
}

// Send는 브랜드메시지를 전송하고 응답을 반환합니다.
//
// Targeting이 M/N/I이면 Contacts가 필요하고 응답 data에 발송 건수(sentCount)가
// 담깁니다. F는 동보 발송이라 Contacts 없이 접수 여부(accepted)만 반환되므로,
// 그 경우 Broadcast를 쓰는 편이 의도가 분명합니다.
func (s *BrandMessageService) Send(req BrandMessageRequest) (map[string]any, error) {
	if req.ScheduleType == "" {
		req.ScheduleType = "DIRECTLY"
	}
	if req.Targeting == "" {
		req.Targeting = "M"
	}
	if req.MessageType == "" {
		req.MessageType = "FT"
	}
	if req.AdFlag == "" {
		req.AdFlag = "Y"
	}
	if req.Adult == "" {
		req.Adult = "N"
	}
	if req.PushAlarm == "" {
		req.PushAlarm = "Y"
	}
	if req.ReplaceSms == "" {
		req.ReplaceSms = "N"
	}
	if req.Buttons == nil {
		req.Buttons = []any{}
	}
	if req.Webhooks == nil {
		req.Webhooks = []string{}
	}

	// 동보는 수신자 목록이 없다. Contacts에 omitempty가 붙어 있어 nil이면
	// 요청 본문에서 아예 빠지고, 빈 배열로 거절되는 일을 막는다.
	if req.Targeting == "F" {
		req.Contacts = nil
	}

	req.KakaoSenderKey = s.kakaoSenderKey
	req.SenderKey = s.smsSenderKey

	return s.http.post("brand-messages/send", req)
}

// Broadcast는 수신 동의한 전체 채널 친구에게 동보 발송합니다(Targeting "F").
//
// 결과는 즉시 알 수 없으므로 Campaigns/Campaign으로 확인합니다.
func (s *BrandMessageService) Broadcast(req BrandMessageRequest) (map[string]any, error) {
	req.Targeting = "F"
	req.Contacts = nil

	return s.Send(req)
}

// Campaigns는 브랜드메시지 캠페인 목록을 조회합니다.
func (s *BrandMessageService) Campaigns(query BrandMessageListQuery) (map[string]any, error) {
	params := map[string]string{}
	if query.From != "" {
		params["from"] = query.From
	}
	if query.To != "" {
		params["to"] = query.To
	}
	if query.Count > 0 {
		params["count"] = strconv.Itoa(query.Count)
	}

	return s.http.get("brand-messages", params)
}

// Campaign은 브랜드메시지 캠페인 상세를 조회합니다.
// campaignID는 발송 응답의 campaignId(UUID)입니다.
func (s *BrandMessageService) Campaign(campaignID string) (map[string]any, error) {
	return s.http.get("brand-messages/"+campaignID, nil)
}
