package sendgo

// SMSService는 SMS/LMS/MMS 전송 서비스입니다.
type SMSService struct {
	http         *httpClient
	smsSenderKey string
}

func newSMSService(hc *httpClient, smsKey string) *SMSService {
	return &SMSService{http: hc, smsSenderKey: smsKey}
}

// SendSMS는 단문 문자(90자 이하)를 전송합니다.
func (s *SMSService) SendSMS(req SmsRequest) error {
	req.MessageType = "SMS"
	return s.send(req)
}

// SendLMS는 장문 문자(2,000자 이하)를 전송합니다.
func (s *SMSService) SendLMS(req SmsRequest) error {
	req.MessageType = "LMS"
	return s.send(req)
}

// SendMMS는 멀티미디어 문자를 전송합니다.
func (s *SMSService) SendMMS(req SmsRequest) error {
	req.MessageType = "MMS"
	return s.send(req)
}

func (s *SMSService) send(req SmsRequest) error {
	if req.CampaignType == "" {
		req.CampaignType = "MESSAGE"
	}
	if req.ScheduleType == "" {
		req.ScheduleType = "DIRECTLY"
	}
	if req.Files == nil {
		req.Files = []any{}
	}
	req.SenderKey = s.smsSenderKey

	_, err := s.http.post("messages/send", req)
	return err
}
