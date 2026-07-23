# sendgo-go

> **Sendgo** Go SDK — 카카오 알림톡/친구톡, SMS/LMS/MMS
> stdlib만 사용 (외부 의존성 없음), Go 1.22+

[![Go Reference](https://pkg.go.dev/badge/github.com/sendgo-dev/sendgo-go.svg)](https://pkg.go.dev/github.com/sendgo-dev/sendgo-go)
[![Go](https://img.shields.io/badge/Go-1.22+-blue)](https://go.dev)

---

## 빠른 시작 (3단계)

### 1단계 — 설치

```bash
go get github.com/sendgo-dev/sendgo-go
```

### 2단계 — 환경변수 설정

```env
SENDGO_ACCESS_KEY=your_access_key
SENDGO_SECRET_KEY=your_secret_key
SENDGO_KAKAO_SENDER_KEY=your_kakao_key
SENDGO_SMS_SENDER_KEY=your_sms_key
```

### 3단계 — 알림톡 전송

```go
package main

import (
    "log"
    "os"

    "github.com/sendgo-dev/sendgo-go/sendgo"
)

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

    err = client.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "ORDER_CONFIRM_001",
        Contacts: []sendgo.Contact{
            {Contact: "01012345678", Name: "홍길동", Var1: "ORD-001"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## 기능별 사용법

### 알림톡

```go
// SMS 대체 발송
subject := "[배송 안내]"
content := "상품이 출고되었습니다."
err = client.Alimtalk.Send(sendgo.AlimtalkRequest{
    TemplateCode: "DELIVERY_001",
    ReplaceSms:   "Y",
    SmsSubject:   &subject,
    SmsContent:   &content,
    Contacts:     []sendgo.Contact{{Contact: "01012345678", Var1: "ORD-001"}},
})

// 예약 발송
at := "2026-04-01 09:00:00"
err = client.Alimtalk.Send(sendgo.AlimtalkRequest{
    TemplateCode: "PROMO_001",
    ScheduleType: "SCHEDULED",
    At:           &at,
    Contacts:     []sendgo.Contact{{Contact: "01012345678"}},
})
```

### SMS / LMS / MMS

```go
// SMS
err = client.SMS.SendSMS(sendgo.SmsRequest{
    Content:  "인증번호: 123456",
    Contacts: []sendgo.Contact{{Contact: "01012345678"}},
})

// LMS
subject := "[공지사항]"
err = client.SMS.SendLMS(sendgo.SmsRequest{
    Subject:  &subject,
    Content:  "서비스 점검이 예정되어 있습니다.",
    Contacts: []sendgo.Contact{{Contact: "01012345678"}},
})
```

### 친구톡

```go
err = client.Friendtalk.Send(sendgo.FriendtalkRequest{
    Content:  "안녕하세요! 이번 주 특가 이벤트입니다.",
    Contacts: []sendgo.Contact{{Contact: "01012345678"}},
})
```

---

## 에러 처리

```go
import "errors"

err = client.Alimtalk.Send(...)
var sendgoErr *sendgo.SendgoError
if errors.As(err, &sendgoErr) {
    fmt.Printf("에러: status=%d code=%s\n", sendgoErr.StatusCode, sendgoErr.ErrorCode)
    switch sendgoErr.ErrorCode {
    case "INVALID_TEMPLATE_CODE":
        fmt.Println("템플릿 코드를 확인하세요.")
    case "PAYMENT_REQUIRED":
        fmt.Println("크레딧이 부족합니다.")
    }
}
```

---

## 라이선스

MIT License © [Sendgo](https://sendgo.io)
