package sendgo

// FriendtalkService는 카카오 친구톡 전송 서비스입니다.
//
// Deprecated: 친구톡은 카카오 정책에 따라 2025-12-31 종료되었습니다.
// 2026-01-01 부터 친구톡 발송 요청은 카카오 측에서 브랜드메시지(자유형)로 자동
// 대체 발송되므로, 이 서비스를 호출해도 실제로 나가는 것은 브랜드메시지입니다.
// 신규 연동은 BrandMessageService를 사용하세요. 다만 자유 본문 타입(FT/FI/FW)을
// 개별 수신자에게 보내는 경로는 아직 이 서비스뿐입니다 — 브랜드메시지 API는 그
// 조합에 NOT_A_BRAND_MESSAGE를 반환합니다. 메시지 타입은 1:1 대응됩니다 —
// FT→BT, FI→BI, FW→BW, FL→BL, FC→BC, FM→BM, FP→BP, FA→BA.
type FriendtalkService struct {
	http           *httpClient
	kakaoSenderKey string
	smsSenderKey   string
}

func newFriendtalkService(hc *httpClient, kakaoKey, smsKey string) *FriendtalkService {
	return &FriendtalkService{http: hc, kakaoSenderKey: kakaoKey, smsSenderKey: smsKey}
}

// Send는 친구톡을 전송합니다.
//
// Deprecated: 2025-12-31 종료. BrandMessageService.Send를 사용하세요.
func (s *FriendtalkService) Send(req FriendtalkRequest) error {
	if req.ScheduleType == "" {
		req.ScheduleType = "DIRECTLY"
	}
	if req.MessageType == "" {
		req.MessageType = "FT"
	}
	if req.AdFlag == "" {
		req.AdFlag = "Y"
	}
	if req.Wide == "" {
		req.Wide = "N"
	}
	if req.Adult == "" {
		req.Adult = "N"
	}
	if req.ReplaceSms == "" {
		req.ReplaceSms = "N"
	}
	if req.Buttons == nil {
		req.Buttons = []any{}
	}
	req.KakaoSenderKey = s.kakaoSenderKey
	req.SenderKey = s.smsSenderKey

	_, err := s.http.post("friends/send", req)
	return err
}
