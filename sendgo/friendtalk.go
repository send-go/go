package sendgo

// FriendtalkService는 카카오 친구톡 전송 서비스입니다.
type FriendtalkService struct {
	http           *httpClient
	kakaoSenderKey string
	smsSenderKey   string
}

func newFriendtalkService(hc *httpClient, kakaoKey, smsKey string) *FriendtalkService {
	return &FriendtalkService{http: hc, kakaoSenderKey: kakaoKey, smsSenderKey: smsKey}
}

// Send는 친구톡을 전송합니다.
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
