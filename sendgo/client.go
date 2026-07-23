package sendgo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type httpClient struct {
	tokenMgr   *tokenManager
	apiVersion string
	baseURL    string
	hc         *http.Client
}

func newHTTPClient(tm *tokenManager, baseURL, apiVersion string) *httpClient {
	return &httpClient{
		tokenMgr:   tm,
		apiVersion: apiVersion,
		baseURL:    baseURL,
		hc:         &http.Client{Timeout: 15 * 1e9},
	}
}

func (c *httpClient) post(path string, body any) (map[string]any, error) {
	return c.doPost(path, body, false)
}

func (c *httpClient) doPost(path string, body any, isRetry bool) (map[string]any, error) {
	token, err := c.tokenMgr.getToken()
	if err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/api/%s/%s", c.baseURL, c.apiVersion, path)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.bearerAuth(token))

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sendgo: request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var responseBody map[string]any
	_ = json.Unmarshal(raw, &responseBody)
	if responseBody == nil {
		responseBody = map[string]any{}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorCode, _ := responseBody["code"].(string)
		parts := strings.Split(path, "/")
		endpoint := parts[len(parts)-1]

		if !isRetry && c.tokenMgr.shouldRefresh(resp.StatusCode, errorCode) {
			c.tokenMgr.invalidate()
			return c.doPost(path, body, true)
		}
		return nil, newSendgoError(resp.StatusCode, responseBody, endpoint, c.apiVersion)
	}

	return responseBody, nil
}

func (c *httpClient) bearerAuth(token string) string {
	if c.apiVersion == "v2" {
		return "Bearer " + token
	}
	return "Bearer " + base64.StdEncoding.EncodeToString([]byte(token))
}
