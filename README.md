# sendgo-go

> **Go에서 카카오 알림톡, 브랜드메시지, SMS를 가장 쉽게 발송하는 공식 Go SDK**

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

> ⚠️ **Deprecated — 친구톡은 카카오 정책에 따라 2025-12-31 종료되었습니다.**
> 2026-01-01 부터 친구톡 발송 요청은 카카오 측에서 **브랜드메시지(자유형)** 로 자동 대체 발송됩니다.
> 호출은 계속 성공하며, 자유 본문 타입(`FT`/`FI`/`FW`)을 개별 수신자에게 보내는 경로는
> 현재 이것뿐이므로 기존 코드를 당장 바꿀 필요는 없습니다.
>
> 다음의 경우에는 **브랜드메시지**를 사용하세요.
> - 템플릿 기반 리치 타입 (`FL`/`FC`/`FM`/`FP`/`FA`)
> - 채널 친구가 **아닌** 수신자 (`targeting` = `N` / `I`)
> - 수신 동의한 전체 채널 친구 동보 (`targeting` = `F`)
>
> 메시지 타입은 1:1 대응되며 변환은 서버가 처리합니다 — `FT`→`BT`, `FI`→`BI`, `FW`→`BW`,
> `FL`→`BL`, `FC`→`BC`, `FM`→`BM`, `FP`→`BP`, `FA`→`BA`.

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

> v2 전용입니다. 자유 본문 타입(`FT`/`FI`/`FW`)을 개별 수신자에게 보낼 때는 여전히 친구톡 API 를 쓰세요 — 이 엔드포인트는 그 조합에 `NOT_A_BRAND_MESSAGE` 를 반환합니다. 친구톡 요청은 카카오 측에서 브랜드메시지(자유형)로 대체 발송됩니다.

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

## 관리 API — 채널·템플릿·발신번호 등록 (v2 전용)

발송은 처음부터 API였지만 **등록과 심사는 콘솔에서만** 되던 것들이 있었습니다.
1.3.0 부터 그 작업도 코드로 처리합니다.

| 서비스 | 하는 일 | 계정 |
| --- | --- | --- |
| `client.KakaoSenders` | 카카오 채널 인증·등록·동기화, 브랜드메시지 M/N 신청 | 기업 |
| `client.NoticeTemplates` | 알림톡 템플릿 CRUD, 검수 요청·취소, 승인 취소, 휴면 해제 | 기업 |
| `client.BrandTemplates` | 브랜드메시지(구 친구톡) 템플릿 CRUD, 동기화, 가져오기 | 기업 |
| `client.SenderRegistration` | 발신번호 등록 신청, 중복 확인, 유형 안내 | 개인·기업 |
| `client.MessageTemplates` | 문자 상용구 템플릿 CRUD | 개인·기업 |
| `client.KakaoImages` | 카카오 이미지 업로드 — 템플릿용 URL 발급 | 기업 |
| `client.RejectedNumbers` | 수신거부(080) 번호 조회 | 개인·기업 |
| `client.Webhook` | 이벤트 웹훅 구독 — 심사 결과 수신 | 개인·기업 |

> **sendgo.io 콘솔에 들어올 일이 없습니다.** 고객의 채널·발신번호·템플릿을
> 여러분 화면만으로 끝까지 처리할 수 있습니다. 휴대폰 발신번호는 콘솔의 PASS
> 본인인증 대신 **신분증 사본(`identityDocument`)을 받아 sendgo 운영자가 대신
> 심사**합니다.
>
> 사람이 개입하는 지점은 **카카오 채널 인증번호 하나**뿐이고, 그마저도
> 여러분 화면에서 입력받으면 됩니다 — 카카오가 관리자 휴대폰으로 직접 보내는
> 확인이라 없앨 수 없습니다.
>
> 심사가 붙는 것들은 **비동기**입니다. 등록 호출이 성공했다는 건 "접수됐다"는
> 뜻이지 "쓸 수 있다"는 뜻이 아닙니다 — 웹훅을 구독해 결과를 받으세요.

### 카카오 채널 등록

```go
// 1단계 — 카카오가 관리자 휴대폰으로 인증번호를 SMS 발송한다 (응답에 번호는 없다)
if _, err := client.KakaoSenders.RequestToken("@my-channel", "01012345678"); err != nil {
    log.Fatal(err)
}

// 2단계 — 사람이 받은 인증번호로 발신프로필 생성
created, err := client.KakaoSenders.Create(sendgo.KakaoSenderCreateRequest{
    Token:        "123456",
    YellowID:     "@my-channel",
    PhoneNumber:  "01012345678",
    CategoryCode: "001001", // Categories("") 로 조회
})
if err != nil {
    log.Fatal(err)
}

sender := created["data"].(map[string]any)["sender"].(map[string]any)
kakaoSenderKey := sender["kakaoSenderKey"].(string)

client.KakaoSenders.Categories("")
client.KakaoSenders.List()
client.KakaoSenders.Sync("")               // 전체 상태 동기화 (하루 한 번 권장)
client.KakaoSenders.Sync(kakaoSenderKey)   // 단건
```

채널이 카카오 쪽에서 차단되면 발송이 조용히 실패하기 시작합니다. `Sync("")` 를
주기적으로 돌리고 `block: true` 인 채널을 감시하세요.

### 알림톡 템플릿 등록과 검수

```go
created, err := client.NoticeTemplates.Create(sendgo.NoticeTemplateRequest{
    KakaoSenderKey:        kakaoSenderKey,
    TemplateName:          "주문 접수 안내",
    TemplateContent:       "#{name}님, 주문 #{orderNo}이 접수되었습니다.",
    TemplateMessageType:   "BA",   // BA 기본형 / EX 부가정보형 / AD 채널추가형 / MI 복합형
    TemplateEmphasizeType: "NONE", // NONE / TEXT / ITEM_LIST / IMAGE
    CategoryCode:          "001001",

    // sendgo 자체 정책 게이트 — 카카오 심사와 별개다
    MessagePurpose:       "order_delivery",
    LegalBasis:           "transaction",
    BenefitOrigin:        "none",
    ExpiryType:           "none",
    OptInReviewConfirmed: true,
    CtaClearConfirmed:    true,
    PolicyConfirmed:      true,
})

template := created["data"].(map[string]any)["template"].(map[string]any)
templateCode := template["templateCode"].(string)

// 검수 요청 — 증빙이 필요하면 파일도 붙인다 (첨부가 있으면 comment 필수)
client.NoticeTemplates.RequestInspection(templateCode, "", nil)

// 결과는 비동기다. 웹훅이 없으므로 폴링한다
synced, _ := client.NoticeTemplates.Sync(templateCode)
status := synced["data"].(map[string]any)["template"].(map[string]any)["inspectionStatus"]
// REG → REQ → APR / REJ
```

정책 필드 조합이 본문과 어긋나면 저장 단계에서 `POLICY_VALIDATION_FAILED` 로
막힙니다. 오류의 `Errors` 에 사유가 한국어로 담기니 그대로 사용자에게 보여
주면 됩니다. 여기서 걸리는 문안은 **카카오 심사에서도 거의 반려**되므로,
며칠 기다렸다 반려당하는 것보다 즉시 아는 편이 낫습니다.

```go
client.NoticeTemplates.List(sendgo.NoticeTemplateListQuery{
    KakaoSenderKey:   kakaoSenderKey,
    InspectionStatus: "APR",
})
client.NoticeTemplates.Show(templateCode)
client.NoticeTemplates.Update(templateCode, req)   // 본문이 바뀌면 재검수 필요
client.NoticeTemplates.CancelInspection(templateCode)
client.NoticeTemplates.CancelApproval(templateCode)
client.NoticeTemplates.Release(templateCode)       // 휴면 해제
client.NoticeTemplates.Delete(templateCode)        // sendgo 목록에서만 삭제된다
client.NoticeTemplates.Categories("")
```

이미지 템플릿과 검수 첨부는 multipart 로 나갑니다.

```go
f, _ := os.Open("banner.jpg")
defer f.Close()

client.NoticeTemplates.CreateWithImage(req, sendgo.MultipartFile{
    FileName: "banner.jpg",
    Content:  f,
})
```

> **삭제 동작이 채널마다 다릅니다.** 알림톡 템플릿은 카카오에 삭제 API 가 없어
> sendgo 목록에서만 빠지고 동기화하면 되살아납니다. 브랜드메시지 템플릿은
> 카카오 쪽에서도 실제로 삭제됩니다.

### 브랜드메시지 템플릿

```go
created, err := client.BrandTemplates.Create(sendgo.BrandTemplateRequest{
    KakaoSenderKey:  kakaoSenderKey,
    TemplateName:    "여름 세일 안내",
    TemplateType:    "FI", // FT/FI/FW/FL/FC/FM/FP/FA — 서버가 chatBubbleType 으로 변환
    TemplateContent: "여름 세일이 시작되었습니다.",
    ImageURL:        "https://mud-kage.kakao.com/....jpg",
})

// 동보 발송(targeting="F")에는 변수가 없는 템플릿만 쓸 수 있다
// created["data"]["template"]["containsVariables"]

client.BrandTemplates.List(sendgo.BrandTemplateListQuery{KakaoSenderKey: kakaoSenderKey})
client.BrandTemplates.Sync(templateCode)
client.BrandTemplates.Import(kakaoSenderKey)   // 카카오에 있는 템플릿 가져오기
client.BrandTemplates.Delete(templateCode)     // 카카오에서도 삭제된다
```

### 발신번호 등록 신청

```go
// 계정 종류에 맞는 유형과 유형별 필수 서류
client.SenderRegistration.NumberTypes()

// 형식·중복 미리 확인
check, _ := client.SenderRegistration.Validate("02-1234-5678", "team_main")

csu, _ := os.Open("csu.pdf")
defer csu.Close()

created, err := client.SenderRegistration.Create(
    sendgo.SenderRegistrationRequest{
        SenderAlias:      "고객센터 대표번호",
        SenderNumberType: "team_main", // personal_other / team_main / team_other_company
        PhoneE164:        "02-1234-5678",
        // check 의 duplicationReasonRequired 가 true 면 필수
        // DuplicationReason: "부서별 분리 운영",
    },
    []sendgo.MultipartFile{
        {FieldName: "csuCertificate", FileName: "csu.pdf", Content: csu},
    },
)

// created["data"]["sender"]["status"] == "PENDING" — 운영자 승인 후 SUCCESS

client.SenderRegistration.List()
client.SenderRegistration.Update(senderKey, "새 이름", "")
client.SenderRegistration.Delete(senderKey)
```

**휴대폰 유형도 API 로 접수할 수 있습니다.** 콘솔의 PASS 본인인증 대신
신분증 사본(`identityDocument`)을 첨부하면 sendgo 운영자가 직접 확인합니다.
이 경로로 접수된 건은 응답의 `identityVerificationMethod` 가 `document` 이고
**자동 승인되지 않습니다** — 운영자 확인 전까지 `PENDING` 입니다.

유형별 필수 서류는 `numberTypes()` 응답의 `requiredDocuments` 로 확인하세요.
반려되면 `rejectionReason` 에 사유가 담깁니다.

### 문자 템플릿

```go
client.MessageTemplates.Create(sendgo.MessageTemplateRequest{
    MessageTranType:    "LMS",
    MessageTranSubject: "주문 안내", // LMS·MMS 는 필수
    MessageTranMsg:     "주문이 접수되었습니다.",
})

client.MessageTemplates.List(sendgo.MessageTemplateListQuery{MessageType: "LMS"})
client.MessageTemplates.Update(templateKey, req)
client.MessageTemplates.Delete(templateKey)
```

### 이벤트 웹훅 — 심사 결과를 밀어 받기

```go
created, _ := client.Webhook.Subscribe(sendgo.WebhookSubscriptionRequest{
    URL:     "https://reseller.example.com/hooks/sendgo",
    Enabled: true,
})

// 시크릿은 이 응답에서 한 번만 나온다. 즉시 저장한다.
secret := created["data"].(map[string]any)["secret"]

client.Webhook.Show()          // 구독 설정 + 마지막 전송 결과
client.Webhook.Test()          // 배선 확인
client.Webhook.Unsubscribe()
```

받는 쪽에서는 **원본 바이트**로 서명을 검증합니다.

```go
func handleWebhook(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)

    if !sendgo.VerifyWebhookSignature(body, r.Header.Get("X-Sendgo-Signature"), secret) {
        http.Error(w, "invalid signature", http.StatusUnauthorized)
        return
    }

    var payload struct {
        Event      string         `json:"event"`
        DeliveryID string         `json:"deliveryId"`
        Data       map[string]any `json:"data"`
    }
    _ = json.Unmarshal(body, &payload)

    // 처리는 큐로. 여기서 오래 끌면 재시도가 쌓인다.
    w.WriteHeader(http.StatusNoContent)
}
```

이벤트 목록은 `sendgo.WebhookEvents` 로 확인할 수 있습니다.

### 카카오 이미지 업로드

브랜드메시지 템플릿의 `imageUrl` 은 **카카오가 호스팅하는 URL** 이어야 합니다.

```go
f, _ := os.Open("banner.jpg")
defer f.Close()

uploaded, _ := client.KakaoImages.Upload("default", sendgo.MultipartFile{
    FileName: "banner.jpg",
    Content:  f,
})

imageURL := uploaded["data"].(map[string]any)["imageUrl"].(string)

client.KakaoImages.UploadMany("carousel_feed", slides)
client.KakaoImages.Types()   // 유형별 필드·최대 개수
```

### 수신거부(080) 동기화

```go
// 증분만 가져간다. 하루 한 번이면 충분하다.
client.RejectedNumbers.List(sendgo.RejectedNumberListQuery{Since: "2026-09-01", Count: 500})
```

---

## 변경 사항

### 1.3.0 (2026-09-11)

- **관리 API 추가** — 콘솔에서만 되던 등록·심사를 코드로 처리합니다.
  `client.KakaoSenders`(채널 인증·등록·동기화, 브랜드메시지 M/N 신청),
  `client.NoticeTemplates`(알림톡 템플릿 CRUD·검수 요청·승인 취소·휴면 해제),
  `client.BrandTemplates`(브랜드메시지 템플릿 CRUD·동기화·가져오기),
  `client.SenderRegistration`(발신번호 등록 신청·중복 확인·유형 안내),
  `client.MessageTemplates`(문자 상용구 템플릿 CRUD).
- `httpClient` 에 `put`·`patch`·`postMultipart` 를 추가했습니다.
  서류 첨부와 이미지 템플릿은 JSON 으로 보낼 수 없습니다. `MultipartFile` 은
  토큰 갱신 재시도를 위해 내용을 미리 버퍼링합니다 — `io.Reader` 는 한 번
  소진되면 되감을 수 없어, 그러지 않으면 재시도가 빈 파일을 올립니다.
- **휴대폰 발신번호도 API 로 접수됩니다.** 콘솔의 PASS 본인인증 대신
  `identityDocument`(신분증 사본)를 첨부하면 sendgo 운영자가 확인합니다.
  이 경로는 자동 승인되지 않고 항상 `PENDING` 으로 시작합니다.
- **이벤트 웹훅** 추가 — 발신번호 승인, 알림톡 검수 결과, 채널 차단,
  브랜드메시지 타겟팅 결과를 구독해 받습니다. 서명은 받은 원본 바이트로
  검증합니다(SDK 에 검증 헬퍼 포함).
- **카카오 이미지 업로드** 추가 — 브랜드메시지 템플릿의 `imageUrl` 은 카카오가
  호스팅하는 URL 이어야 하는데, 그 URL 을 얻는 길이 콘솔에만 있었습니다.
- **수신거부(080) 조회** 추가 — 자기 DB 의 수신 상태를 맞출 수 있습니다.

### 1.2.0 (2026-08-14)

- **친구톡 Deprecated 표기** — 친구톡은 카카오 정책에 따라 2025-12-31 종료되었고,
  2026-01-01 부터 발송 요청이 브랜드메시지(자유형)로 자동 대체 발송됩니다.
  관련 API 에 각 언어의 표준 deprecation 표기를 달았습니다.
- 자유 본문 타입(`FT`/`FI`/`FW`)의 개별 발송 경로는 아직 친구톡 API 뿐이라는 점을
  문서에 명시했습니다 — 브랜드메시지 API 는 그 조합에 `NOT_A_BRAND_MESSAGE` 를 반환합니다.
- 브랜드메시지 전환 안내와 메시지 타입 1:1 대응표를 README 에 추가했습니다.

### 1.1.0 (2026-08-11)

- 짧은 URL 추가 — `client.ShortURL`
- **모듈 경로 수정** — `github.com/sendgo-dev/sendgo-go` → `github.com/send-go/go`. README·문서가 모두 후자를 쓰고 있어 `go get` 이 동작할 수 없었다.
- `sendgo.String()` 헬퍼 추가. `At`/`SmsSubject`/`SmsContent`/`Subject` 가 `*string` 이라 리터럴을 직접 넘길 수 없었다.
- `httpClient.delete()` 추가

## 라이선스

MIT License © 2026 [Sendgo](https://sendgo.io)

---

*키워드: 카카오 알림톡 Go, 카카오 친구톡 Golang, SMS 발송 Go, 알림톡 Go SDK, Go 카카오 API 연동, Sendgo Go SDK, Golang 알림 발송*

## 계정·조직·API 키 관리 (1.5.0)

발송용 `accessKey`/`secretKey`가 없는 단계에서 사용하는 **별도 계정 클라이언트**입니다.
콘솔에서 발급받은 에이전트 토큰(`SENDGO_AGENT_TOKEN`)으로 `/api/v2/account`를 호출합니다.
계정 조회에는 `account:read`, 키·허용 IP 변경에는 `keys:write` 권한이 필요합니다.
토큰 만료나 권한 부족(401/403)은 그대로 예외로 반환하며 자동 갱신·재시도하지 않습니다.

조직 선택은 서버에 저장되는 **사용자 계정의 현재 조직**을 바꿉니다. 같은 사용자로
여러 조직의 설정을 동시에 변경하지 마세요. 개인 계정으로 돌아가려면 조직 ID에
`null`(Python `None`, Ruby `nil`, Go `nil`) 또는 `personal`을 전달합니다.
키 발급 응답의 `data.apiKey.secretKey`는 한 번만 반환되므로 서버의 비밀 저장소에 보관하세요.
허용 IP가 하나라도 등록되면 목록 밖의 IP는 차단됩니다.
에이전트 토큰과 키는 브라우저·모바일 앱에 포함하거나 응답·로그에 출력하지 않습니다.

```go
account, err := sendgo.NewAccount(os.Getenv("SENDGO_AGENT_TOKEN"), "")
if err != nil { log.Fatal(err) }
result, err := account.Me()
if err != nil { log.Fatal(err) }
_ = result
```

지원 메서드: `Me`, `Organizations`, `SelectOrganization`, `ApiKeys`, `CreateApiKey`, `ApiKey`, `UpdateApiKey`, `DeleteApiKey`, `IssueToken`, `AllowedIps`, `AddAllowedIp`, `DeleteAllowedIp`.

키 생성 인자는 `name`, 선택적 `ipAddresses: [{ip, description}]`이며, 허용 IP 추가 인자는 `ip`, 선택적 `description`입니다. 키·IP 식별자는 응답의 `id`(UUID)를 사용합니다.

## 템플릿 폴더 (1.5.0)

기업 계정의 발송용 API 키와 `apiVersion=v2` 설정으로 사용하는 서버 전용 API입니다.
폴더는 알림톡·브랜드메시지가 공유하며, 목록의 `templateType`은 `notice` 또는 `brand`입니다.
목록은 `data.folders` 트리와 `total`, `uncategorised` 개수를 반환합니다.
`templateCount`는 하위 폴더를 제외한 해당 폴더의 템플릿 수입니다.

- 생성: `name`, 선택 `parentUuid`. 최대 5단계이며 같은 부모 아래 이름 중복은 409입니다.
- 이동: 동일 발신프로필의 `templateCodes` 1~100개. `folderUuid`는 필수이며 `null`이면 미분류로 이동합니다.
- 템플릿 목록: `folderUuid=none`은 미분류, UUID는 해당 폴더, 생략은 전체입니다.
- 템플릿 등록: 선택 필드 `folderUuid`로 폴더를 지정합니다. 기존 템플릿 수정 API 대신 폴더 이동 API를 사용하세요.

승인되지 않은 키의 `403 ACCESS_KEY_NOT_APPROVED`는 토큰 재발급·재시도 없이 반환합니다.
계정 API의 `autoApprove`는 서버 설정의 실제 승인 정책을 나타냅니다.

```go
folders, err := client.TemplateFolders.List(sendgo.TemplateFolderListQuery{TemplateType: "notice"})
if err != nil { return err }
_ = folders
_, err = client.TemplateFolders.Create(sendgo.TemplateFolderCreateRequest{Name: "주문 안내"})
if err != nil { return err }
_, err = client.TemplateFolders.Assign(sendgo.TemplateFolderAssignRequest{
    TemplateType: "notice", KakaoSenderKey: kakaoSenderKey,
    TemplateCodes: []string{"ORDER_001"}, FolderUUID: nil,
})
if err != nil { return err }
_, err = client.NoticeTemplates.List(sendgo.NoticeTemplateListQuery{FolderUUID: "none"})
// 템플릿 생성 요청의 FolderUUID 필드로 폴더를 지정합니다.
```
