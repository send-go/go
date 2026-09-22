package sendgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// AccountClient는 서버 전용 계정 API 클라이언트입니다. 에이전트 토큰은 자동 갱신하지 않습니다.
type AccountClient struct {
	agentToken string
	baseURL    string
	http       *http.Client
}

// NewAccount는 발송용 키 없이 에이전트 토큰으로 계정 클라이언트를 만듭니다.
// baseURL이 비어 있으면 https://sendgo.io를 사용합니다.
func NewAccount(agentToken, baseURL string) (*AccountClient, error) {
	if strings.TrimSpace(agentToken) == "" {
		return nil, fmt.Errorf("sendgo: agentToken은 필수입니다")
	}
	if baseURL == "" {
		baseURL = "https://sendgo.io"
	}
	return &AccountClient{agentToken: agentToken, baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{Timeout: 15 * time.Second, CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}}, nil
}

// Me: 계정 상태와 다음 단계 조회.
func (c *AccountClient) Me() (map[string]any, error) {
	return c.request("GET", "", nil)
}

// Organizations: 조직 목록 조회.
func (c *AccountClient) Organizations() (map[string]any, error) {
	return c.request("GET", "organizations", nil)
}

// SelectOrganization: 조직 선택. null은 개인 계정.
func (c *AccountClient) SelectOrganization(organizationId *string) (map[string]any, error) {
	return c.request("POST", "organizations/select", map[string]any{"organizationId": organizationId})
}

// ApiKeys: 현재 조직의 API 키 목록.
func (c *AccountClient) ApiKeys() (map[string]any, error) {
	return c.request("GET", "api-keys", nil)
}

// CreateApiKey: API 키 발급. secretKey는 이 응답에서만 반환.
func (c *AccountClient) CreateApiKey(params map[string]any) (map[string]any, error) {
	return c.request("POST", "api-keys", params)
}

// ApiKey: API 키 상세 조회.
func (c *AccountClient) ApiKey(apiKeyId string) (map[string]any, error) {
	return c.request("GET", "api-keys/"+url.PathEscape(apiKeyId)+"", nil)
}

// UpdateApiKey: API 키 이름 변경.
func (c *AccountClient) UpdateApiKey(apiKeyId string, name string) (map[string]any, error) {
	return c.request("PATCH", "api-keys/"+url.PathEscape(apiKeyId)+"", map[string]any{"name": name})
}

// DeleteApiKey: API 키 폐기.
func (c *AccountClient) DeleteApiKey(apiKeyId string) (map[string]any, error) {
	return c.request("DELETE", "api-keys/"+url.PathEscape(apiKeyId)+"", nil)
}

// IssueToken: 승인된 API 키의 발송용 토큰 발급.
func (c *AccountClient) IssueToken(apiKeyId string) (map[string]any, error) {
	return c.request("POST", "api-keys/"+url.PathEscape(apiKeyId)+"/token", map[string]any{})
}

// AllowedIps: 허용 IP 목록과 호출자 IP 조회.
func (c *AccountClient) AllowedIps(apiKeyId string) (map[string]any, error) {
	return c.request("GET", "api-keys/"+url.PathEscape(apiKeyId)+"/allowed-ips", nil)
}

// AddAllowedIp: 허용 IP 추가. ip와 선택적 description 사용.
func (c *AccountClient) AddAllowedIp(apiKeyId string, params map[string]any) (map[string]any, error) {
	return c.request("POST", "api-keys/"+url.PathEscape(apiKeyId)+"/allowed-ips", params)
}

// DeleteAllowedIp: 허용 IP 삭제.
func (c *AccountClient) DeleteAllowedIp(apiKeyId string, ipId string) (map[string]any, error) {
	return c.request("DELETE", "api-keys/"+url.PathEscape(apiKeyId)+"/allowed-ips/"+url.PathEscape(ipId)+"", nil)
}

func (c *AccountClient) request(method, path string, body any) (map[string]any, error) {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(payload)
	}
	endpoint := c.baseURL + "/api/v2/account"
	if path != "" {
		endpoint += "/" + path
	}
	req, err := http.NewRequest(method, endpoint, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.agentToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := map[string]any{}
	decodeErr := json.Unmarshal(raw, &data)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, newSendgoError(resp.StatusCode, data, path, "v2")
	}
	if len(raw) > 0 && decodeErr != nil {
		return nil, decodeErr
	}
	return data, nil
}
