# sendgo-go

> **Go에서 카카오 알림톡, 친구톡, SMS를 가장 쉽게 발송하는 공식 Go SDK**

[![Go Reference](https://pkg.go.dev/badge/github.com/send-go/go.svg)](https://pkg.go.dev/github.com/send-go/go)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

`sendgo-go`는 [Sendgo](https://sendgo.io) 알림 API를 위한 공식 Go SDK입니다.
**표준 라이브러리만 사용**하며, 완전한 타입 안전성과 동시성 안전 토큰 관리를 제공합니다.

---

## 설치

```bash
go get github.com/send-go/go
```

---

## 빠른 시작

```go
package main

import (
    "fmt"
    "log"
    "os"

    sendgo "github.com/send-go/go/sendgo"
)

func main() {
    client, err := sendgo.New(sendgo.Config{
        AccessKey:      os.Getenv("SENDGO_ACCESS_KEY"),
        SecretKey:      os.Getenv("SENDGO_SECRET_KEY"),
        KakaoSenderKey: os.Getenv("SENDGO_KAKAO_SENDER_KEY"),
        SmsSenderKey:   os.Getenv("SENDGO_SMS_SENDER_KEY"),
        ApiVersion:     "v2",
    })
    if err != nil {
        log.Fatal(err)
    }

    // 알림톡 발송
    result, err := client.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "ORDER_CONFIRM_001",
        Contacts: []sendgo.Contact{
            {Contact: "01012345678", Name: "홍길동", Var1: "ORD-001", Var2: "29,000원"},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("발송 결과: %v\n", result)
}
```

---

## 알림톡 상세 사용법

```go
package main

import (
    "log"
    "os"
    sendgo "github.com/send-go/go/sendgo"
)

func main() {
    client, _ := sendgo.New(sendgo.Config{
        AccessKey:      os.Getenv("SENDGO_ACCESS_KEY"),
        SecretKey:      os.Getenv("SENDGO_SECRET_KEY"),
        KakaoSenderKey: os.Getenv("SENDGO_KAKAO_SENDER_KEY"),
        SmsSenderKey:   os.Getenv("SENDGO_SMS_SENDER_KEY"),
        ApiVersion:     "v2",
    })

    // 다건 발송
    client.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "ORDER_CONFIRM_001",
        Contacts: []sendgo.Contact{
            {Contact: "01011111111", Name: "홍길동", Var1: "ORD-001", Var2: "29,000원"},
            {Contact: "01022222222", Name: "김철수", Var1: "ORD-002", Var2: "15,000원"},
            {Contact: "01033333333", Name: "이영희", Var1: "ORD-003", Var2: "52,000원"},
        },
    })

    // 예약 발송
    client.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "PROMO_SUMMER_2026",
        ScheduleType: "SCHEDULED",
        At:           "2026-07-28 09:00:00",
        Contacts: []sendgo.Contact{
            {Contact: "01012345678", Var1: "여름 한정 50% 할인"},
        },
    })

    // SMS 자동 대체 발송
    client.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "DELIVERY_START_001",
        ReplaceSms:   "Y",
        SmsSubject:   "[배송 시작 안내]",
        SmsContent:   "주문하신 상품이 출고되었습니다.\n송장번호: #{var2}",
        Contacts: []sendgo.Contact{
            {Contact: "01012345678", Var1: "ORD-001", Var2: "1234567890"},
        },
    })
}
```

---

## 친구톡 사용법

```go
// 텍스트형
client.Friendtalk.Send(sendgo.FriendtalkRequest{
    Content: "안녕하세요! 7월 한정 특가 이벤트를 확인해보세요.",
    Contacts: []sendgo.Contact{
        {Contact: "01012345678"},
    },
})

// 이미지형
client.Friendtalk.Send(sendgo.FriendtalkRequest{
    MessageType: "FI",
    Content:     "이번 주 특가 상품을 확인하세요!",
    ImageURL:    "https://cdn.example.com/banner.jpg",
    ImageLink:   "https://example.com/event",
    Contacts: []sendgo.Contact{
        {Contact: "01012345678"},
    },
})
```

---

## SMS / LMS / MMS 사용법

```go
// SMS
client.SMS.SendSMS(sendgo.SmsRequest{
    Content: "[Sendgo] 인증번호: 123456 (5분 이내 입력)",
    Contacts: []sendgo.Contact{
        {Contact: "01012345678"},
    },
})

// LMS
client.SMS.SendLMS(sendgo.SmsRequest{
    Subject: "[중요] 서비스 점검 안내",
    Content: "안녕하세요. 서비스 점검이 예정되어 있습니다.\n■ 일시: 2026-07-25 02:00 ~ 06:00",
    Contacts: []sendgo.Contact{
        {Contact: "01012345678"},
    },
})

// MMS
client.SMS.SendMMS(sendgo.SmsRequest{
    Subject: "[이벤트] 7월 특가",
    Content: "이번 달 특가 상품을 확인하세요!",
    Contacts: []sendgo.Contact{
        {Contact: "01011111111"},
        {Contact: "01022222222"},
    },
})

// 예약 문자
client.SMS.SendSMS(sendgo.SmsRequest{
    Content:      "[알림] 예약 미팅을 확인해주세요.",
    ScheduleType: "SCHEDULED",
    At:           "2026-07-23 08:00:00",
    Contacts: []sendgo.Contact{{Contact: "01012345678"}},
})
```

---

## 프레임워크 통합

### net/http 서버

```go
package main

import (
    "encoding/json"
    "net/http"
    "os"
    sendgo "github.com/send-go/go/sendgo"
)

var sg *sendgo.Client

func init() {
    var err error
    sg, err = sendgo.New(sendgo.Config{
        AccessKey:      os.Getenv("SENDGO_ACCESS_KEY"),
        SecretKey:      os.Getenv("SENDGO_SECRET_KEY"),
        KakaoSenderKey: os.Getenv("SENDGO_KAKAO_SENDER_KEY"),
        ApiVersion:     "v2",
    })
    if err != nil {
        panic(err)
    }
}

func notifyOrderHandler(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Phone   string `json:"phone"`
        OrderNo string `json:"order_no"`
    }
    json.NewDecoder(r.Body).Decode(&req)

    _, err := sg.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "ORDER_CONFIRM_001",
        Contacts:     []sendgo.Contact{{Contact: req.Phone, Var1: req.OrderNo}},
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func main() {
    http.HandleFunc("/api/notify/order", notifyOrderHandler)
    http.ListenAndServe(":8080", nil)
}
```

### Gin Framework

```go
package main

import (
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    sendgo "github.com/send-go/go/sendgo"
)

func main() {
    sg, _ := sendgo.New(sendgo.Config{
        AccessKey:      os.Getenv("SENDGO_ACCESS_KEY"),
        SecretKey:      os.Getenv("SENDGO_SECRET_KEY"),
        KakaoSenderKey: os.Getenv("SENDGO_KAKAO_SENDER_KEY"),
        ApiVersion:     "v2",
    })

    r := gin.Default()

    r.POST("/api/notify/order", func(c *gin.Context) {
        var req struct {
            Phone   string `json:"phone"`
            OrderNo string `json:"order_no"`
        }
        c.ShouldBindJSON(&req)

        _, err := sg.Alimtalk.Send(sendgo.AlimtalkRequest{
            TemplateCode: "ORDER_CONFIRM_001",
            Contacts:     []sendgo.Contact{{Contact: req.Phone, Var1: req.OrderNo}},
        })
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{"success": true})
    })

    r.Run(":8080")
}
```

---

## 예외 처리

```go
import "github.com/send-go/go/sendgo"

_, err := client.Alimtalk.Send(req)
if err != nil {
    var se *sendgo.SendgoError
    if errors.As(err, &se) {
        fmt.Printf("발송 실패: HTTP %d [%s]\n", se.StatusCode, se.ErrorCode)
        switch se.ErrorCode {
        case "INVALID_ACCESS_KEY", "INVALID_SECRET_KEY":
            alertOps("Sendgo 인증키를 확인하세요.")
        case "INVALID_TEMPLATE_CODE":
            log.Printf("존재하지 않는 템플릿: %s", se.Message)
        case "PAYMENT_REQUIRED":
            alertOps("Sendgo 크레딧이 부족합니다.")
        }
    }
}
```

---

## 설정 옵션

| 필드 | 타입 | 필수 | 기본값 | 설명 |
|------|------|------|--------|------|
| `AccessKey` | `string` | **필수** | — | Sendgo 액세스 키 |
| `SecretKey` | `string` | **필수** | — | Sendgo 시크릿 키 |
| `KakaoSenderKey` | `string` | 선택 | `""` | 카카오 발신프로필 키 |
| `SmsSenderKey` | `string` | 선택 | `""` | SMS 발신자 키 |
| `ApiVersion` | `string` | 선택 | `"v2"` | API 버전 (`v1` \| `v2`) |
| `BaseURL` | `string` | 선택 | `"https://sendgo.io"` | API 기본 URL |

---

## 관련 패키지

| 언어/프레임워크 | 패키지 | GitHub |
|----------------|--------|--------|
| Spring Boot | `io.sendgo:sendgo-spring` | [spring](https://github.com/send-go/spring) |
| Node.js | `@sendgo/node` | [node](https://github.com/send-go/node) |
| Python | `sendgo-python` | [python](https://github.com/send-go/python) |
| PHP | `sendgo/php` | [php](https://github.com/send-go/php) |
| 전체 목록 | — | [send-go GitHub 조직](https://github.com/send-go) |

---

## 라이선스

MIT License © 2026 [Sendgo](https://sendgo.io)

---

*키워드: 카카오 알림톡 Go, 카카오 친구톡 Golang, SMS 발송 Go, 알림톡 Go SDK, Go 카카오 API 연동, Sendgo Go SDK, Golang 알림 발송*
