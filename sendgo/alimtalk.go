package sendgo

// AlimtalkService는 카카오 알림톡 전송 서비스입니다.
type AlimtalkService struct {
	http           *httpClient
	kakaoSenderKey string
	smsSenderKey   string
}

func newAlimtalkService(hc *httpClient, kakaoKey, smsKey string) *AlimtalkService {
	return &AlimtalkService{http: hc, kakaoSenderKey: kakaoKey, smsSenderKey: smsKey}
}

// Send는 알림톡을 전송합니다.
// req.ScheduleType 기본값: "DIRECTLY", req.ReplaceSms 기본값: "N"
func (s *AlimtalkService) Send(req AlimtalkRequest) error {
	if req.ScheduleType == "" {
		req.ScheduleType = "DIRECTLY"
	}
	if req.ReplaceSms == "" {
		req.ReplaceSms = "N"
	}
	req.KakaoSenderKey = s.kakaoSenderKey
	req.SenderKey = s.smsSenderKey

	_, err := s.http.post("notices/send", req)
	return err
}
