package sendgo

// String은 문자열 리터럴을 *string 필드에 바로 넘길 수 있게 해주는 헬퍼입니다.
//
// AlimtalkRequest.At / SmsSubject / SmsContent 와 SmsRequest.Subject 는
// "설정하지 않음"(JSON null)과 "빈 문자열"을 구분해야 하므로 *string 입니다.
// Go에서는 리터럴의 주소를 직접 얻을 수 없어(&"..." 는 컴파일 에러) 임시 변수를
// 만들어야 했는데, 이 헬퍼를 쓰면 호출부가 한 줄로 정리됩니다.
//
//	err := client.SMS.SendLMS(sendgo.SmsRequest{
//	    Subject:  sendgo.String("[중요] 서비스 점검 안내"),
//	    Content:  "...",
//	    Contacts: []sendgo.Contact{{Contact: "01012345678"}},
//	})
func String(s string) *string {
	return &s
}
