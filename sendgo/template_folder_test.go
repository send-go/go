package sendgo

import (
	"errors"
	"os"
	"testing"
)

// TestTemplateFolderContract는 모의 HTTP 서버와 공개 클라이언트의 계약을 검증합니다.
func TestTemplateFolderContract(t *testing.T) {
	base := os.Getenv("SENDGO_TEST_URL")
	if base == "" {
		t.Skip("모의 HTTP 서버가 필요합니다")
	}
	c, err := New(Config{AccessKey: "test-access", SecretKey: "test-secret", APIVersion: "v2", BaseURL: base})
	if err != nil {
		t.Fatal(err)
	}
	check := func(_ map[string]any, err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	f := "11111111-1111-4111-8111-111111111111"
	key := "채널 /?"
	check(c.TemplateFolders.List(TemplateFolderListQuery{}))
	check(c.TemplateFolders.List(TemplateFolderListQuery{TemplateType: "brand", KakaoSenderKey: key}))
	check(c.TemplateFolders.Create(TemplateFolderCreateRequest{Name: "주문"}))
	check(c.TemplateFolders.Create(TemplateFolderCreateRequest{Name: "하위", ParentUUID: &f}))
	for _, kind := range []string{"notice", "brand"} {
		check(c.TemplateFolders.Assign(TemplateFolderAssignRequest{TemplateType: kind, KakaoSenderKey: key, TemplateCodes: []string{"코드 1", "code/2"}, FolderUUID: &f}))
		check(c.TemplateFolders.Assign(TemplateFolderAssignRequest{TemplateType: kind, KakaoSenderKey: key, TemplateCodes: []string{"코드 1"}, FolderUUID: nil}))
	}
	check(c.NoticeTemplates.List(NoticeTemplateListQuery{FolderUUID: "none"}))
	check(c.BrandTemplates.List(BrandTemplateListQuery{FolderUUID: f}))
	check(c.NoticeTemplates.Create(NoticeTemplateRequest{TemplateName: "테스트", FolderUUID: f}))
	check(c.BrandTemplates.Create(BrandTemplateRequest{TemplateName: "테스트", FolderUUID: f}))
	for _, tc := range []struct {
		kind   string
		status int
		code   string
	}{{"forbidden", 403, "ACCESS_KEY_NOT_APPROVED"}, {"invalid", 422, "VALIDATION_FAILED"}, {"missing", 404, "TEMPLATE_FOLDER_NOT_FOUND"}} {
		_, err := c.TemplateFolders.List(TemplateFolderListQuery{TemplateType: tc.kind})
		var apiErr *SendgoError
		if !errors.As(err, &apiErr) || apiErr.StatusCode != tc.status || apiErr.ErrorCode != tc.code {
			t.Fatalf("예상하지 못한 오류: %v", err)
		}
	}
}
