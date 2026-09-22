package sendgo

import (
	"errors"
	"os"
	"testing"
)

func TestAccountContract(t *testing.T) {
	base := os.Getenv("SENDGO_TEST_URL")
	if base == "" {
		t.Skip("_sdktest/account-contract.py로 모의 서버와 함께 실행")
	}
	if _, err := NewAccount("", base); err == nil {
		t.Fatal("빈 토큰 허용")
	}
	c, err := NewAccount("test-agent", base+"/")
	if err != nil {
		t.Fatal(err)
	}
	check := func(data map[string]any, err error) {
		t.Helper()
		if err != nil || data["message"] != "Success" {
			t.Fatalf("응답 오류: %v %v", data, err)
		}
	}
	organization := "team-id"
	check(c.Me())
	check(c.Organizations())
	check(c.SelectOrganization(nil))
	check(c.SelectOrganization(&organization))
	check(c.ApiKeys())
	check(c.CreateApiKey(map[string]any{"name": "한글 이름", "ipAddresses": []any{map[string]any{"ip": "192.0.2.1", "description": "서버"}}}))
	check(c.ApiKey("key/id ?"))
	check(c.UpdateApiKey("key/id ?", "새 이름"))
	check(c.DeleteApiKey("key/id ?"))
	check(c.IssueToken("key/id ?"))
	check(c.AllowedIps("key/id ?"))
	check(c.AddAllowedIp("key/id ?", map[string]any{"ip": "192.0.2.1", "description": "서버"}))
	check(c.DeleteAllowedIp("key/id ?", "ip/id ?"))
	for _, tt := range []struct {
		token  string
		status int
		code   string
	}{{"expired", 401, "AGENT_TOKEN_EXPIRED"}, {"forbidden", 403, "AGENT_ABILITY_MISSING"}} {
		client, _ := NewAccount(tt.token, base)
		_, err := client.Me()
		var apiError *SendgoError
		if !errors.As(err, &apiError) || apiError.StatusCode != tt.status || apiError.ErrorCode != tt.code {
			t.Fatalf("잘못된 오류: %v", err)
		}
	}
}
