package sendgo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var noRefreshCodes = map[string]bool{
	"INVALID_AUTH_HEADER": true, "INVALID_BASIC_AUTH": true, "INVALID_BASIC_AUTH_PAYLOAD": true,
	"INVALID_ACCESS_KEY": true, "INVALID_SECRET_KEY": true, "ACCESS_KEY_NOT_APPROVED": true,
	"TEAM_REQUIRED_FOR_KAKAO": true, "IP_NOT_ALLOWED": true,
	"INVALID_SENDER_KEY": true, "INVALID_KAKAO_SENDER_KEY": true,
}

const tokenTTL = 50 * time.Minute

// tokenManager는 Sendgo API 토큰을 인메모리 캐시로 관리합니다.
type tokenManager struct {
	baseURL    string
	accessKey  string
	secretKey  string
	apiVersion string
	mu         sync.Mutex
	token      string
	expiresAt  time.Time
}

func newTokenManager(baseURL, accessKey, secretKey, apiVersion string) *tokenManager {
	return &tokenManager{
		baseURL:    baseURL,
		accessKey:  accessKey,
		secretKey:  secretKey,
		apiVersion: apiVersion,
	}
}

func (tm *tokenManager) getToken() (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.token != "" && time.Now().Before(tm.expiresAt) {
		return tm.token, nil
	}
	return tm.fetchToken()
}

func (tm *tokenManager) invalidate() {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.token = ""
	tm.expiresAt = time.Time{}
}

func (tm *tokenManager) shouldRefresh(status int, errorCode string) bool {
	if status != 401 && status != 403 {
		return false
	}
	if tm.apiVersion == "v2" && noRefreshCodes[errorCode] {
		return false
	}
	return true
}

func (tm *tokenManager) fetchToken() (string, error) {
	url := fmt.Sprintf("%s/api/%s/token", tm.baseURL, tm.apiVersion)
	creds := base64.StdEncoding.EncodeToString([]byte(tm.accessKey + ":" + tm.secretKey))

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+creds)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("sendgo: token request failed: %w", err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		body = map[string]any{}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", newSendgoError(resp.StatusCode, body, "token", tm.apiVersion)
	}

	data, _ := body["data"].(map[string]any)
	token, _ := data["token"].(string)
	if token == "" {
		return "", fmt.Errorf("sendgo: token field missing in response")
	}

	tm.token = token
	tm.expiresAt = time.Now().Add(tokenTTL)
	return token, nil
}
