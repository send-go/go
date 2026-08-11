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
        At:           sendgo.String("2026-07-28 09:00:00"),
        Contacts: []sendgo.Contact{
            {Contact: "01012345678", Var1: "여름 한정 50% 할인"},
        },
    })

    // SMS 자동 대체 발송
    client.Alimtalk.Send(sendgo.AlimtalkRequest{
        TemplateCode: "DELIVERY_START_001",
        ReplaceSms:   "Y",
        SmsSubject:   sendgo.String("[배송 시작 안내]"),
        SmsContent:   sendgo.String("주문하신 상품이 출고되었습니다.\n송장번호: #{var2}"),
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

## 브랜드메시지 사용법

브랜드메시지는 친구톡의 후속 채널입니다. 메시지 타입이 친구톡과 1:1 대응되며
(`FT`→`BT`, `FI`→`BI`, `FW`→`BW`, `FL`→`BL`, `FC`→`BC`, `FM`→`BM`, `FP`→`BP`, `FA`→`BA`),
요청에는 **친구톡 코드를 그대로** 넘기고 변환은 서버가 처리합니다.

친구톡과 달리 다음이 가능합니다.

- 채널 친구가 **아닌** 수신자에게 발송 (`targeting: N`)
- 수신 동의한 **전체 채널 친구 동보** 발송 (`targeting: F`, 수신자 목록 불필요)
- 리스트·캐러셀·커머스·동영상 등 **템플릿 기반 리치 메시지**

> v2 전용입니다. `FT`/`FI`/`FW`를 채널 친구에게만 보낼 때는 친구톡 API가 더 간단합니다.

```go
// 단건 발송 — 채널 친구 대상
_, err := client.BrandMessage.Send(sendgo.BrandMessageRequest{
    Targeting:          "M",
    MessageType:        "FL",
    FriendTemplateUUID: "9cd5460b-6458-4edc-9b11-c26d3013c340",
    Contacts: []sendgo.Contact{
        {Contact: "01012345678", Var1: "29,000원"},
    },
})

// 동보 발송 — 수신 동의한 전체 채널 친구 (Contacts 불필요)
_, err = client.BrandMessage.Broadcast(sendgo.BrandMessageRequest{
    MessageType:        "FW",
    FriendTemplateUUID: "9cd5460b-6458-4edc-9b11-c26d3013c340",
})

// 캠페인 조회
list, err := client.BrandMessage.Campaigns(sendgo.BrandMessageListQuery{Count: 10})
one, err := client.BrandMessage.Campaign("1f0a6d0e-6b3b-4f0f-9b2f-2f6f6a1b7c11")
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
    Subject: sendgo.String("[중요] 서비스 점검 안내"),
    Content: "안녕하세요. 서비스 점검이 예정되어 있습니다.\n■ 일시: 2026-07-25 02:00 ~ 06:00",
    Contacts: []sendgo.Contact{
        {Contact: "01012345678"},
    },
})

// MMS
client.SMS.SendMMS(sendgo.SmsRequest{
    Subject: sendgo.String("[이벤트] 7월 특가"),
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
    At:           sendgo.String("2026-07-23 08:00:00"),
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

## 포인터 필드와 `sendgo.String`

`AlimtalkRequest` 의 `At`/`SmsSubject`/`SmsContent` 와 `SmsRequest` 의 `Subject` 는
"설정하지 않음"(JSON `null`)과 빈 문자열을 구분해야 하므로 `*string` 입니다.
Go에서는 리터럴의 주소를 얻을 수 없으므로(`&"..."` 는 컴파일 에러) 변수를 따로
선언하거나 `sendgo.String(...)` 헬퍼를 사용하세요.

```go
err := client.SMS.SendLMS(sendgo.SmsRequest{
    Subject:  sendgo.String("[중요] 서비스 점검 안내"),
    Content:  "...",
    Contacts: []sendgo.Contact{{Contact: "01012345678"}},
})
```

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

## 짧은 URL

짧은 URL 은 메시지 본문의 링크를 줄이고, 그 링크가 실제로 눌렸는지 집계합니다.
문자는 바이트 수가 요금과 직결되므로 링크를 줄이면 그만큼 본문을 더 쓸 수 있습니다.

같은 원본 URL 을 다시 줄이면 **기존 링크가 그대로 반환**됩니다. 캠페인별로 반응을
따로 집계하려면 `forceNew` 로 새 코드를 만드세요.

`deactivate` 는 링크를 삭제하지 않고 리다이렉트만 중지합니다. 이미 발송한 메시지의
링크를 무효화할 때 쓰며, 누적 통계는 남고 이후 접속은 `410 Gone` 이 됩니다.

```go
// 짧은 URL 생성 (v2 전용)
created, err := client.ShortURL.Create(sendgo.ShortURLRequest{
	TargetURL: "https://example.com/promotions/summer-sale",
	Title:     "여름 세일 랜딩",
})
if err != nil {
	log.Fatal(err)
}

code := created["data"].(map[string]any)["code"].(string)

// 반응 통계 — 일별 추이 + 디바이스/유입경로/국가별 분해
stats, err := client.ShortURL.Stats(code, sendgo.ShortURLStatsQuery{From: "2026-08-01"})

client.ShortURL.List(sendgo.ShortURLListQuery{Count: 10})
client.ShortURL.Show(code)
client.ShortURL.Deactivate(code) // 리다이렉트만 중지, 통계는 남는다
```

`stats` 는 일별 추이(`daily`)와 디바이스(`byDevice`)·유입경로(`byReferer`)·국가(`byCountry`)별
분해를 반환합니다. 일별 추이는 사전 집계 표에서 읽으므로 클릭이 많아도 응답 시간이 일정합니다.

## 변경 사항

### 1.1.0 (2026-08-11)

- 짧은 URL 추가 — `client.ShortURL`
- **모듈 경로 수정** — `github.com/sendgo-dev/sendgo-go` → `github.com/send-go/go`. README·문서가 모두 후자를 쓰고 있어 `go get` 이 동작할 수 없었다.
- `sendgo.String()` 헬퍼 추가. `At`/`SmsSubject`/`SmsContent`/`Subject` 가 `*string` 이라 리터럴을 직접 넘길 수 없었다.
- `httpClient.delete()` 추가

## 라이선스

MIT License © 2026 [Sendgo](https://sendgo.io)

---

*키워드: 카카오 알림톡 Go, 카카오 친구톡 Golang, SMS 발송 Go, 알림톡 Go SDK, Go 카카오 API 연동, Sendgo Go SDK, Golang 알림 발송*
