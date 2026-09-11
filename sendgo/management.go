package sendgo

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
)

// structToFields는 구조체를 multipart 필드 맵으로 바꿉니다.
//
// multipart 로 보낼 때도 JSON 태그의 이름을 그대로 써야 서버가 같은 필드로
// 읽습니다. json.Marshal 을 한 번 거치면 omitempty 처리와 태그 이름 매핑을
// 손으로 다시 적지 않아도 됩니다.
func structToFields(v any) (map[string]any, error) {
	encoded, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("sendgo: failed to encode request: %w", err)
	}

	var fields map[string]any
	if err := json.Unmarshal(encoded, &fields); err != nil {
		return nil, fmt.Errorf("sendgo: failed to decode request: %w", err)
	}

	return fields, nil
}

// ---------------------------------------------------------------- 요청 타입

// KakaoSenderCreateRequest는 카카오 발신프로필 등록 요청(2단계)입니다.
type KakaoSenderCreateRequest struct {
	// Token은 1단계에서 관리자 휴대폰으로 받은 인증번호입니다.
	Token string `json:"token"`
	// YellowID는 채널 검색용 아이디입니다. "@"는 있어도 없어도 됩니다.
	YellowID string `json:"yellowId"`
	// PhoneNumber는 채널 관리자 휴대폰 번호입니다.
	PhoneNumber string `json:"phoneNumber"`
	// CategoryCode는 Categories()로 조회한 코드입니다.
	CategoryCode string `json:"categoryCode"`
}

// NoticeTemplateRequest는 알림톡 템플릿 등록·수정 요청입니다.
//
// 뒤쪽 정책 필드 일곱 개는 sendgo 자체 게이트입니다. 카카오 심사와 별개이며
// 조합이 본문과 어긋나면 POLICY_VALIDATION_FAILED로 거절됩니다.
type NoticeTemplateRequest struct {
	KakaoSenderKey string `json:"kakaoSenderKey,omitempty"`
	TemplateName   string `json:"templateName"`
	// TemplateContent의 변수는 #{name} 형식으로 씁니다.
	TemplateContent string `json:"templateContent"`
	// TemplateMessageType: BA 기본형 / EX 부가정보형 / AD 채널추가형 / MI 복합형.
	TemplateMessageType string `json:"templateMessageType"`
	// TemplateEmphasizeType: NONE / TEXT / ITEM_LIST / IMAGE.
	TemplateEmphasizeType string `json:"templateEmphasizeType"`
	// CategoryCode는 6자리 숫자입니다.
	CategoryCode string `json:"categoryCode"`

	TemplateTitle         string           `json:"templateTitle,omitempty"`
	TemplateSubtitle      string           `json:"templateSubtitle,omitempty"`
	TemplateHeader        string           `json:"templateHeader,omitempty"`
	TemplateExtra         string           `json:"templateExtra,omitempty"`
	TemplateItem          map[string]any   `json:"templateItem,omitempty"`
	TemplateItemHighlight map[string]any   `json:"templateItemHighlight,omitempty"`
	TemplateRepresentLink map[string]any   `json:"templateRepresentLink,omitempty"`
	Buttons               []map[string]any `json:"buttons,omitempty"`
	QuickReplies          []map[string]any `json:"quickReplies,omitempty"`
	SecurityFlag          bool             `json:"securityFlag,omitempty"`
	AdultFlag             bool             `json:"adultFlag,omitempty"`

	// MessagePurpose: order_delivery / reservation_booking / payment_billing /
	// account_auth / service_ops / policy_notice / benefit_notice /
	// customer_support / other.
	MessagePurpose string `json:"messagePurpose"`
	// LegalBasis: transaction / paid_purchase / event_entry / contract / policy_notice.
	LegalBasis string `json:"legalBasis"`
	// BenefitOrigin: none / paid / event / contract / promo / free.
	BenefitOrigin string `json:"benefitOrigin"`
	// ExpiryType: none / rights_based / promo.
	ExpiryType string `json:"expiryType"`
	// OptInReviewConfirmed는 true여야 합니다.
	OptInReviewConfirmed bool `json:"optInReviewConfirmed"`
	// CtaClearConfirmed는 true여야 합니다.
	CtaClearConfirmed bool `json:"ctaClearConfirmed"`
	// PolicyConfirmed는 true여야 검수를 요청할 수 있습니다.
	PolicyConfirmed bool `json:"policyConfirmed"`
}

// NoticeTemplateListQuery는 알림톡 템플릿 목록 조회 조건입니다.
type NoticeTemplateListQuery struct {
	KakaoSenderKey string
	// InspectionStatus: REG / REQ / APR / REJ / BLOCK / DORMANT.
	InspectionStatus string
	Search           string
	Count            int
}

func (q NoticeTemplateListQuery) toMap() map[string]string {
	m := map[string]string{}
	if q.KakaoSenderKey != "" {
		m["kakaoSenderKey"] = q.KakaoSenderKey
	}
	if q.InspectionStatus != "" {
		m["inspectionStatus"] = q.InspectionStatus
	}
	if q.Search != "" {
		m["search"] = q.Search
	}
	if q.Count > 0 {
		m["count"] = fmt.Sprintf("%d", q.Count)
	}
	return m
}

// BrandTemplateRequest는 브랜드메시지 템플릿 등록·수정 요청입니다.
//
// TemplateType은 친구톡 표기(FT/FI/FW/FL/FC/FM/FP/FA)를 그대로 씁니다 —
// 서버가 chatBubbleType으로 변환합니다.
type BrandTemplateRequest struct {
	KakaoSenderKey    string           `json:"kakaoSenderKey,omitempty"`
	TemplateName      string           `json:"templateName"`
	TemplateType      string           `json:"templateType"`
	TemplateContent   string           `json:"templateContent,omitempty"`
	Adult             bool             `json:"adult,omitempty"`
	Header            string           `json:"header,omitempty"`
	AdditionalContent string           `json:"additional_content,omitempty"`
	ImageURL          string           `json:"imageUrl,omitempty"`
	ImageLink         string           `json:"imageLink,omitempty"`
	Buttons           []map[string]any `json:"buttons,omitempty"`
	Coupon            map[string]any   `json:"coupon,omitempty"`
	Item              map[string]any   `json:"item,omitempty"`
	Commerce          map[string]any   `json:"commerce,omitempty"`
	List              []map[string]any `json:"list,omitempty"`
	Head              map[string]any   `json:"head,omitempty"`
	Tail              map[string]any   `json:"tail,omitempty"`
	Video             map[string]any   `json:"video,omitempty"`
	MainWideItem      map[string]any   `json:"mainWideItem,omitempty"`
	SubWideItemList   []map[string]any `json:"subWideItemList,omitempty"`
}

// BrandTemplateListQuery는 브랜드메시지 템플릿 목록 조회 조건입니다.
type BrandTemplateListQuery struct {
	KakaoSenderKey string
	Search         string
	Count          int
}

func (q BrandTemplateListQuery) toMap() map[string]string {
	m := map[string]string{}
	if q.KakaoSenderKey != "" {
		m["kakaoSenderKey"] = q.KakaoSenderKey
	}
	if q.Search != "" {
		m["search"] = q.Search
	}
	if q.Count > 0 {
		m["count"] = fmt.Sprintf("%d", q.Count)
	}
	return m
}

// SenderRegistrationRequest는 발신번호 등록 신청의 텍스트 필드입니다.
// 서류는 []MultipartFile로 따로 넘깁니다.
type SenderRegistrationRequest struct {
	// SenderAlias는 계정 안에서 중복될 수 없는 관리용 이름입니다.
	SenderAlias string
	// SenderNumberType: 여섯 유형 전부 API로 접수할 수 있습니다.
	// 휴대폰 계열은 PASS 대신 identityDocument(신분증 사본)를 첨부합니다.
	SenderNumberType string
	// PhoneE164는 숫자와 하이픈만 씁니다. 서버가 E.164로 정규화합니다.
	PhoneE164 string
	// DuplicationReason은 Validate()의 duplicationReasonRequired가 true면 필수입니다.
	DuplicationReason string
	// 아래 셋은 team_other_company 필수입니다.
	AcceptanceName   string
	DelegationName   string
	DelegationReason string
}

func (r SenderRegistrationRequest) toFields() map[string]any {
	fields := map[string]any{
		"senderAlias":      r.SenderAlias,
		"senderNumberType": r.SenderNumberType,
		"phoneE164":        r.PhoneE164,
	}
	// 빈 값을 보내면 서버가 "빈 값으로 저장"으로 읽는다. 채워진 것만 넣는다.
	optional := map[string]string{
		"duplicationReason": r.DuplicationReason,
		"acceptanceName":    r.AcceptanceName,
		"delegationName":    r.DelegationName,
		"delegationReason":  r.DelegationReason,
	}
	for key, value := range optional {
		if value != "" {
			fields[key] = value
		}
	}
	return fields
}

// MessageTemplateRequest는 문자 상용구 템플릿 등록·수정 요청입니다.
type MessageTemplateRequest struct {
	// MessageTranType: SMS / LMS / MMS.
	MessageTranType string `json:"messageTranType"`
	MessageTranMsg  string `json:"messageTranMsg"`
	// MessageTranSubject는 LMS·MMS 필수입니다. SMS에 넣으면 발송 시 버려집니다.
	MessageTranSubject string `json:"messageTranSubject,omitempty"`
	IsFavorite         bool   `json:"isFavorite,omitempty"`
}

// MessageTemplateListQuery는 문자 템플릿 목록 조회 조건입니다.
type MessageTemplateListQuery struct {
	// MessageType: SMS / LMS / MMS / TOTAL.
	MessageType string
	Search      string
	Count       int
}

func (q MessageTemplateListQuery) toMap() map[string]string {
	m := map[string]string{}
	if q.MessageType != "" {
		m["messageType"] = q.MessageType
	}
	if q.Search != "" {
		m["search"] = q.Search
	}
	if q.Count > 0 {
		m["count"] = fmt.Sprintf("%d", q.Count)
	}
	return m
}

// ---------------------------------------------------- 카카오 발신프로필(채널)

// KakaoSenderService는 카카오 발신프로필(채널) 관리 서비스입니다.
//
// v2 전용이며 기업(Team) 소유 애플리케이션만 사용할 수 있습니다.
//
// 채널 등록은 두 단계입니다. 카카오가 인증번호를 채널 관리자 휴대폰으로 SMS
// 발송하므로 완전 무인 자동화는 불가능합니다 — 사람이 문자를 받아 Create에
// 넣어야 합니다.
//
// 사용법:
//
//	// 1단계 — 관리자 휴대폰으로 인증번호 발송 (응답에 번호는 없습니다)
//	_, err := client.KakaoSenders.RequestToken("@my-channel", "01012345678")
//
//	// 2단계 — 사람이 받은 인증번호로 발신프로필 생성
//	created, err := client.KakaoSenders.Create(sendgo.KakaoSenderCreateRequest{
//	    Token:        "123456",
//	    YellowID:     "@my-channel",
//	    PhoneNumber:  "01012345678",
//	    CategoryCode: "001001",
//	})
type KakaoSenderService struct {
	http *httpClient
}

func newKakaoSenderService(hc *httpClient) *KakaoSenderService {
	return &KakaoSenderService{http: hc}
}

// RequestToken은 1단계 — 채널 인증번호를 발송합니다.
//
// 응답에 인증번호는 들어있지 않습니다. 카카오가 phoneNumber로 SMS를 보냅니다.
func (s *KakaoSenderService) RequestToken(yellowID, phoneNumber string) (map[string]any, error) {
	return s.http.post("kakao-senders/token", map[string]string{
		"yellowId":    yellowID,
		"phoneNumber": phoneNumber,
	})
}

// Create는 2단계 — 발신프로필을 등록합니다.
//
// 이미 등록된 채널을 다시 등록해도 오류가 아닙니다. 카카오가 같은 senderKey를
// 돌려주고 서버가 기존 행을 갱신합니다.
func (s *KakaoSenderService) Create(req KakaoSenderCreateRequest) (map[string]any, error) {
	return s.http.post("kakao-senders", req)
}

// List는 발신프로필 목록을 조회합니다.
func (s *KakaoSenderService) List() (map[string]any, error) {
	return s.http.get("kakao-senders", nil)
}

// Show는 발신프로필 상세를 조회합니다.
func (s *KakaoSenderService) Show(kakaoSenderKey string) (map[string]any, error) {
	return s.http.get("kakao-senders/"+url.PathEscape(kakaoSenderKey), nil)
}

// Categories는 발신프로필 카테고리를 조회합니다.
// categoryCode가 비어 있으면 전체를 반환합니다.
func (s *KakaoSenderService) Categories(categoryCode string) (map[string]any, error) {
	query := map[string]string{}
	if categoryCode != "" {
		query["categoryCode"] = categoryCode
	}
	return s.http.get("kakao-senders/categories", query)
}

// Sync는 발신프로필 상태를 카카오에서 다시 읽어 옵니다.
// kakaoSenderKey가 비어 있으면 팀 전체를 동기화합니다.
//
// 채널이 카카오 쪽에서 차단·휴면되면 발송이 조용히 실패하기 시작합니다.
// 그 사실을 먼저 알 방법은 이 호출뿐이므로 하루 한 번 정도 돌리는 게 좋습니다.
func (s *KakaoSenderService) Sync(kakaoSenderKey string) (map[string]any, error) {
	if kakaoSenderKey == "" {
		return s.http.post("kakao-senders/sync", map[string]any{})
	}
	return s.http.post("kakao-senders/"+url.PathEscape(kakaoSenderKey)+"/sync", map[string]any{})
}

// UploadBrandMessageEvidence는 브랜드메시지 M 신청에 필요한 광고성 정보
// 수신동의 증적자료를 업로드합니다. jpg/png, 5MB 이하.
func (s *KakaoSenderService) UploadBrandMessageEvidence(kakaoSenderKey string, evidence MultipartFile) (map[string]any, error) {
	evidence.FieldName = "evidence"
	return s.http.postMultipart(
		"kakao-senders/"+url.PathEscape(kakaoSenderKey)+"/brand-message/evidence",
		nil,
		[]MultipartFile{evidence},
	)
}

// ApplyBrandMessageTargeting은 브랜드메시지 M(마케팅) / N(정보성) 사용을
// 신청합니다.
//
// 결과는 즉시 확정되지 않습니다. 발신프로필의 brandMessageStatus로 확인합니다.
func (s *KakaoSenderService) ApplyBrandMessageTargeting(kakaoSenderKey, targetType string) (map[string]any, error) {
	return s.http.post(
		"kakao-senders/"+url.PathEscape(kakaoSenderKey)+"/brand-message/apply",
		map[string]string{"targetType": targetType},
	)
}

// ------------------------------------------------------------ 알림톡 템플릿

// NoticeTemplateService는 알림톡 템플릿 관리 서비스입니다.
//
// v2 전용이며 기업(Team) 소유 애플리케이션만 사용할 수 있습니다.
//
// 템플릿은 만든 즉시 쓸 수 없습니다. 카카오 검수를 통과해야 합니다:
//
//	등록      inspectionStatus=REG   ← 발송 불가
//	검수 요청  inspectionStatus=REQ   ← 카카오 심사 중
//	승인      inspectionStatus=APR   ← 여기부터 발송 가능
//	반려      inspectionStatus=REJ   ← comments 에 사유
//
// 검수 결과는 비동기입니다. 웹훅이 없으므로 Sync로 폴링합니다.
type NoticeTemplateService struct {
	http *httpClient
}

func newNoticeTemplateService(hc *httpClient) *NoticeTemplateService {
	return &NoticeTemplateService{http: hc}
}

// List는 알림톡 템플릿 목록을 조회합니다.
func (s *NoticeTemplateService) List(query NoticeTemplateListQuery) (map[string]any, error) {
	return s.http.get("notice-templates", query.toMap())
}

// Show는 템플릿 상세를 조회합니다.
// 응답의 data.template.policy 에 정책 검토 상태가 들어 있습니다.
func (s *NoticeTemplateService) Show(templateCode string) (map[string]any, error) {
	return s.http.get(noticeTemplatePath(templateCode), nil)
}

// Create는 템플릿을 등록합니다.
func (s *NoticeTemplateService) Create(req NoticeTemplateRequest) (map[string]any, error) {
	return s.http.post("notice-templates", req)
}

// CreateWithImage는 이미지 템플릿을 등록합니다
// (TemplateEmphasizeType이 "IMAGE"인 경우).
//
// multipart로 나가므로 Buttons 같은 필드는 JSON 문자열로 직렬화해 보냅니다.
func (s *NoticeTemplateService) CreateWithImage(req NoticeTemplateRequest, image MultipartFile) (map[string]any, error) {
	image.FieldName = "image"
	fields, err := structToFields(req)
	if err != nil {
		return nil, err
	}
	return s.http.postMultipart("notice-templates", fields, []MultipartFile{image})
}

// Update는 템플릿을 수정합니다.
//
// 발신프로필과 템플릿 코드는 바꿀 수 없습니다. 본문·버튼처럼 카카오에 등록된
// 내용이 바뀌면 검수 상태가 되돌아가므로 재검수를 요청해야 합니다.
func (s *NoticeTemplateService) Update(templateCode string, req NoticeTemplateRequest) (map[string]any, error) {
	return s.http.put(noticeTemplatePath(templateCode), req)
}

// Delete는 템플릿을 삭제합니다.
//
// 카카오는 템플릿 삭제 API를 제공하지 않습니다. sendgo 목록에서만 지워지고
// 비즈니스 채널 쪽 템플릿은 남습니다. 동기화하면 다시 나타납니다.
func (s *NoticeTemplateService) Delete(templateCode string) (map[string]any, error) {
	return s.http.delete(noticeTemplatePath(templateCode))
}

// Sync는 카카오에서 검수 상태와 반려 사유를 다시 읽어 옵니다.
func (s *NoticeTemplateService) Sync(templateCode string) (map[string]any, error) {
	return s.http.post(noticeTemplatePath(templateCode)+"/sync", map[string]any{})
}

// RequestInspection은 검수를 요청합니다.
//
// 첨부가 있으면 comment는 필수입니다. 정책 검토를 통과하지 못한 템플릿은
// POLICY_REVIEW_REQUIRED로 거절되고 errors.reasons에 사유가 담깁니다.
func (s *NoticeTemplateService) RequestInspection(templateCode, comment string, attachments []MultipartFile) (map[string]any, error) {
	path := noticeTemplatePath(templateCode) + "/inspection"

	if len(attachments) == 0 {
		body := map[string]any{}
		if comment != "" {
			body["comment"] = comment
		}
		return s.http.post(path, body)
	}

	// 서버는 attachments[0], attachments[1] 형태를 기대합니다.
	for i := range attachments {
		attachments[i].FieldName = fmt.Sprintf("attachments[%d]", i)
	}

	fields := map[string]any{}
	if comment != "" {
		fields["comment"] = comment
	}

	return s.http.postMultipart(path, fields, attachments)
}

// CancelInspection은 검수 요청을 취소합니다. 아직 심사 중(REQ)일 때만 통합니다.
func (s *NoticeTemplateService) CancelInspection(templateCode string) (map[string]any, error) {
	return s.http.delete(noticeTemplatePath(templateCode) + "/inspection")
}

// CancelApproval은 승인을 취소합니다. 이후에는 발송할 수 없습니다.
func (s *NoticeTemplateService) CancelApproval(templateCode string) (map[string]any, error) {
	return s.http.delete(noticeTemplatePath(templateCode) + "/approval")
}

// Release는 휴면을 해제합니다.
func (s *NoticeTemplateService) Release(templateCode string) (map[string]any, error) {
	return s.http.post(noticeTemplatePath(templateCode)+"/release", map[string]any{})
}

// Categories는 템플릿 카테고리 코드를 조회합니다.
func (s *NoticeTemplateService) Categories(categoryCode string) (map[string]any, error) {
	query := map[string]string{}
	if categoryCode != "" {
		query["categoryCode"] = categoryCode
	}
	return s.http.get("notice-templates/categories", query)
}

func noticeTemplatePath(templateCode string) string {
	return "notice-templates/" + url.PathEscape(templateCode)
}

// ------------------------------------------------------ 브랜드메시지 템플릿

// BrandTemplateService는 브랜드메시지(구 친구톡) 템플릿 관리 서비스입니다.
//
// v2 전용이며 기업(Team) 소유 애플리케이션만 사용할 수 있습니다.
// 알림톡 템플릿과 달리 검수 요청 단계가 없습니다.
type BrandTemplateService struct {
	http *httpClient
}

func newBrandTemplateService(hc *httpClient) *BrandTemplateService {
	return &BrandTemplateService{http: hc}
}

// List는 브랜드메시지 템플릿 목록을 조회합니다.
func (s *BrandTemplateService) List(query BrandTemplateListQuery) (map[string]any, error) {
	return s.http.get("brand-templates", query.toMap())
}

// Show는 템플릿 상세를 조회합니다.
// sendgo 코드(KFT-...)와 카카오 브랜드 템플릿 코드 둘 다 받습니다.
func (s *BrandTemplateService) Show(templateCode string) (map[string]any, error) {
	return s.http.get(brandTemplatePath(templateCode), nil)
}

// Create는 템플릿을 등록합니다.
func (s *BrandTemplateService) Create(req BrandTemplateRequest) (map[string]any, error) {
	return s.http.post("brand-templates", req)
}

// Update는 템플릿을 수정합니다. 발신프로필은 바꿀 수 없습니다.
func (s *BrandTemplateService) Update(templateCode string, req BrandTemplateRequest) (map[string]any, error) {
	return s.http.put(brandTemplatePath(templateCode), req)
}

// Delete는 템플릿을 삭제합니다. 알림톡과 달리 카카오 쪽에서도 실제로 삭제됩니다.
func (s *BrandTemplateService) Delete(templateCode string) (map[string]any, error) {
	return s.http.delete(brandTemplatePath(templateCode))
}

// Sync는 카카오에서 템플릿 상태를 다시 읽어 옵니다.
// 카카오 쪽에서 이미 삭제됐으면 로컬에서도 제거하고 data.deleted: true를 반환합니다.
func (s *BrandTemplateService) Sync(templateCode string) (map[string]any, error) {
	return s.http.post(brandTemplatePath(templateCode)+"/sync", map[string]any{})
}

// Import는 발신프로필 단위로 카카오 쪽에 이미 있는 템플릿을 들여옵니다.
func (s *BrandTemplateService) Import(kakaoSenderKey string) (map[string]any, error) {
	return s.http.post("brand-templates/import", map[string]string{"kakaoSenderKey": kakaoSenderKey})
}

func brandTemplatePath(templateCode string) string {
	return "brand-templates/" + url.PathEscape(templateCode)
}

// ---------------------------------------------------------------- 발신번호

// RegistrableSenderTypes는 API로 접수할 수 있는 발신번호 유형입니다 — 전부입니다.
var RegistrableSenderTypes = []string{
	"personal_mobile",
	"personal_other",
	"team_main",
	"team_representative_mobile",
	"team_emp_mobile",
	"team_other_company",
}

// IdentityDocumentSenderTypes는 신분증 사본(identityDocument)이 필요한 유형입니다.
//
// 콘솔은 PASS 본인인증을 쓰지만 API는 신분증 사본을 받아 sendgo 운영자가
// 직접 확인합니다. 이 경로로 접수된 건은 자동 승인되지 않습니다.
var IdentityDocumentSenderTypes = []string{"personal_mobile", "team_representative_mobile", "team_emp_mobile"}

// SenderRegistrationService는 발신번호 등록·심사 접수 서비스입니다.
//
// v2 전용. 카카오와 달리 개인 계정 애플리케이션도 쓸 수 있습니다.
//
// 등록하면 곧바로 쓸 수 있는 게 아니라 PENDING으로 접수되고, 운영자 승인 후
// SUCCESS가 됩니다.
type SenderRegistrationService struct {
	http *httpClient
}

func newSenderRegistrationService(hc *httpClient) *SenderRegistrationService {
	return &SenderRegistrationService{http: hc}
}

// List는 발신번호 목록을 조회합니다. 심사 상태(status)를 여기서 확인합니다.
func (s *SenderRegistrationService) List() (map[string]any, error) {
	return s.http.get("senders", nil)
}

// Show는 발신번호 상세를 조회합니다.
func (s *SenderRegistrationService) Show(senderKey string) (map[string]any, error) {
	return s.http.get("senders/"+url.PathEscape(senderKey), nil)
}

// NumberTypes는 계정 종류에 맞는 발신번호 유형과 유형별 필수 서류를 반환합니다.
//
// 유형별 identityVerification(none/document)과 필요한 서류 목록을 줍니다.
func (s *SenderRegistrationService) NumberTypes() (map[string]any, error) {
	return s.http.get("senders/number-types", nil)
}

// Validate는 등록 전 형식·중복을 확인합니다.
//
// 응답의 duplicationReasonRequired가 true면 Create에 DuplicationReason을
// 함께 넣어야 합니다.
func (s *SenderRegistrationService) Validate(phoneE164, senderNumberType string) (map[string]any, error) {
	return s.http.post("senders/validate", map[string]string{
		"phoneE164":        phoneE164,
		"senderNumberType": senderNumberType,
	})
}

// Create는 발신번호 등록을 신청합니다. 서류가 붙으므로 multipart로 나갑니다.
//
// files에는 최소한 csuCertificate(통신서비스 이용증명원)가 있어야 합니다.
// 휴대폰 계열은 identityDocument(신분증 사본)가, team_other_company는
// 수임·위임 서류가 더 필요합니다 — NumberTypes()로 확인하세요.
func (s *SenderRegistrationService) Create(req SenderRegistrationRequest, files []MultipartFile) (map[string]any, error) {
	return s.http.postMultipart("senders", req.toFields(), files)
}

// Update는 별칭과 기본 발신 지정을 변경합니다. 번호와 심사 상태는 바꿀 수 없습니다.
func (s *SenderRegistrationService) Update(senderKey, senderAlias, primaryType string) (map[string]any, error) {
	body := map[string]string{"senderAlias": senderAlias}
	if primaryType != "" {
		body["primaryType"] = primaryType
	}
	return s.http.patch("senders/"+url.PathEscape(senderKey), body)
}

// Delete는 발신번호를 삭제합니다.
// 기본 발신번호를 지우면 남은 번호 중 하나가 기본으로 승계됩니다.
func (s *SenderRegistrationService) Delete(senderKey string) (map[string]any, error) {
	return s.http.delete("senders/" + url.PathEscape(senderKey))
}

// ------------------------------------------------------------- 문자 템플릿

// MessageTemplateService는 문자(SMS/LMS/MMS) 상용구 템플릿 서비스입니다.
//
// v2 전용. 카카오 템플릿과 달리 검수가 없어 만들면 바로 쓸 수 있고,
// 기업 계정이 아니어도 됩니다.
type MessageTemplateService struct {
	http *httpClient
}

func newMessageTemplateService(hc *httpClient) *MessageTemplateService {
	return &MessageTemplateService{http: hc}
}

// List는 문자 템플릿 목록을 조회합니다.
func (s *MessageTemplateService) List(query MessageTemplateListQuery) (map[string]any, error) {
	return s.http.get("message-templates", query.toMap())
}

// Show는 문자 템플릿 상세를 조회합니다.
func (s *MessageTemplateService) Show(templateKey string) (map[string]any, error) {
	return s.http.get("message-templates/"+url.PathEscape(templateKey), nil)
}

// Create는 문자 템플릿을 등록합니다. LMS·MMS는 MessageTranSubject가 필수입니다.
func (s *MessageTemplateService) Create(req MessageTemplateRequest) (map[string]any, error) {
	return s.http.post("message-templates", req)
}

// Update는 문자 템플릿을 수정합니다.
func (s *MessageTemplateService) Update(templateKey string, req MessageTemplateRequest) (map[string]any, error) {
	return s.http.put("message-templates/"+url.PathEscape(templateKey), req)
}

// Delete는 문자 템플릿을 삭제합니다 (소프트 삭제 — 목록에서만 사라집니다).
func (s *MessageTemplateService) Delete(templateKey string) (map[string]any, error) {
	return s.http.delete("message-templates/" + url.PathEscape(templateKey))
}

// ---------------------------------------------------------------- 웹훅

// WebhookEvents 는 구독할 수 있는 이벤트 목록입니다.
var WebhookEvents = []string{
	"sender.status_changed",
	"notice_template.inspection_status_changed",
	"kakao_sender.status_changed",
	"kakao_sender.brand_message_status_changed",
}

// WebhookSubscriptionRequest 는 웹훅 구독 생성·수정 요청입니다.
type WebhookSubscriptionRequest struct {
	// URL 은 이벤트를 받을 주소입니다. https 만 허용됩니다.
	URL string `json:"url"`
	// Secret 은 서명 키입니다. 비우면 서버가 만들어 응답에서 한 번만 돌려줍니다.
	// 이미 있는 상태에서 비우면 기존 값을 유지합니다.
	Secret string `json:"secret,omitempty"`
	// Events 를 비우면 전체 구독입니다.
	Events []string `json:"events,omitempty"`
	// Enabled 는 구독 활성 여부입니다.
	Enabled bool `json:"enabled"`
}

// WebhookService 는 이벤트 웹훅 구독 서비스입니다.
//
// v2 전용. 심사는 비동기라 폴링 말고는 방법이 없었습니다. 구독해 두면 상태가
// 바뀔 때마다 도착합니다.
type WebhookService struct {
	http *httpClient
}

func newWebhookService(hc *httpClient) *WebhookService {
	return &WebhookService{http: hc}
}

// Show 는 현재 구독 설정을 조회합니다.
// 마지막 전송 결과(latStatus)도 함께 오므로 내 엔드포인트가 실제로 받고
// 있는지 확인할 수 있습니다.
func (s *WebhookService) Show() (map[string]any, error) {
	return s.http.get("webhook", nil)
}

// Subscribe 는 구독을 만들거나 수정합니다.
func (s *WebhookService) Subscribe(req WebhookSubscriptionRequest) (map[string]any, error) {
	return s.http.put("webhook", req)
}

// Test 는 테스트 이벤트를 보냅니다. 구독 목록과 무관하게 도착합니다.
func (s *WebhookService) Test() (map[string]any, error) {
	return s.http.post("webhook/test", map[string]any{})
}

// Unsubscribe 는 구독을 해지합니다.
func (s *WebhookService) Unsubscribe() (map[string]any, error) {
	return s.http.delete("webhook")
}

// VerifyWebhookSignature 는 수신한 웹훅의 서명을 검증합니다.
//
// rawBody 는 **받은 바이트 그대로**여야 합니다. 파싱한 뒤 다시 인코딩한 값으로
// 계산하면 키 순서나 이스케이프 차이로 검증이 깨집니다.
//
//	body, _ := io.ReadAll(r.Body)
//	if !sendgo.VerifyWebhookSignature(body, r.Header.Get("X-Sendgo-Signature"), secret) {
//	    http.Error(w, "invalid signature", http.StatusUnauthorized)
//	    return
//	}
func VerifyWebhookSignature(rawBody []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(rawBody)
	expected := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(signature))
}

// ------------------------------------------------------------ 카카오 이미지

// SingleKakaoImageTypes 는 파일 하나를 올리고 URL 하나를 받는 유형입니다.
var SingleKakaoImageTypes = []string{
	"alimtalk",
	"alimtalk_highlight",
	"default",
	"wide",
	"wide_item_list_first",
}

// MultiKakaoImageTypes 는 파일 여러 개를 올리는 유형과 최대 개수입니다.
var MultiKakaoImageTypes = map[string]int{
	"wide_item_list":    4,
	"carousel_feed":     10,
	"carousel_commerce": 11,
}

// KakaoImageService 는 카카오 이미지 업로드 서비스입니다.
//
// v2 전용, 기업 계정 전용. 브랜드메시지 템플릿의 imageUrl 은 아무 URL 이나
// 되는 게 아니라 **카카오가 호스팅하는 URL** 이어야 하고, 그 URL 을 얻는
// 방법이 이 업로드뿐입니다.
type KakaoImageService struct {
	http *httpClient
}

func newKakaoImageService(hc *httpClient) *KakaoImageService {
	return &KakaoImageService{http: hc}
}

// Types 는 업로드 가능한 유형과 제약을 조회합니다.
func (s *KakaoImageService) Types() (map[string]any, error) {
	return s.http.get("kakao-images/types", nil)
}

// Upload 는 단일 이미지를 올리고 data.imageUrl 을 받습니다. jpg/png, 2MB 이하.
func (s *KakaoImageService) Upload(imageType string, image MultipartFile) (map[string]any, error) {
	image.FieldName = "image"

	return s.http.postMultipart("kakao-images/"+url.PathEscape(imageType), nil, []MultipartFile{image})
}

// UploadMany 는 다중 이미지를 올립니다. 유형별 최대 개수가 다릅니다.
func (s *KakaoImageService) UploadMany(imageType string, images []MultipartFile) (map[string]any, error) {
	// 서버는 images[0], images[1] 형태를 기대합니다.
	named := make([]MultipartFile, len(images))
	for i, image := range images {
		image.FieldName = fmt.Sprintf("images[%d]", i)
		named[i] = image
	}

	return s.http.postMultipart("kakao-images/"+url.PathEscape(imageType), nil, named)
}

// ---------------------------------------------------------------- 수신거부

// RejectedNumberListQuery 는 수신거부 번호 조회 조건입니다.
type RejectedNumberListQuery struct {
	// Since 이후 수신거부된 건만. 증분 동기화에 씁니다.
	Since  string
	Search string
	Count  int
}

func (q RejectedNumberListQuery) toMap() map[string]string {
	m := map[string]string{}
	if q.Since != "" {
		m["since"] = q.Since
	}
	if q.Search != "" {
		m["search"] = q.Search
	}
	if q.Count > 0 {
		m["count"] = fmt.Sprintf("%d", q.Count)
	}
	return m
}

// RejectedNumberService 는 수신거부(080) 번호 조회 서비스입니다. 조회 전용.
//
// 발송 API 가 알아서 제외하지만 **자기 DB 의 수신 상태도 맞춰야** 합니다 —
// 그러지 않으면 매번 보내고 매번 걸러지는 것을 반복하고, 자기 화면에서는
// 여전히 "수신 동의" 로 보입니다.
type RejectedNumberService struct {
	http *httpClient
}

func newRejectedNumberService(hc *httpClient) *RejectedNumberService {
	return &RejectedNumberService{http: hc}
}

// List 는 수신거부 번호 목록을 조회합니다.
func (s *RejectedNumberService) List(query RejectedNumberListQuery) (map[string]any, error) {
	return s.http.get("rejected-numbers", query.toMap())
}
