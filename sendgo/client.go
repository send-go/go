package sendgo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
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
	return c.do(http.MethodPost, path, body, nil, false)
}

// get performs a GET request, optionally with a query string. Used by the
// campaign lookup endpoints, which have no request body.
func (c *httpClient) get(path string, query map[string]string) (map[string]any, error) {
	return c.do(http.MethodGet, path, nil, query, false)
}

// delete performs a DELETE request. `do` already drives the verb, so this only
// needs to omit the body.
func (c *httpClient) delete(path string) (map[string]any, error) {
	return c.do(http.MethodDelete, path, nil, nil, false)
}

func (c *httpClient) do(method, path string, body any, query map[string]string, isRetry bool) (map[string]any, error) {
	token, err := c.tokenMgr.getToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/%s/%s", c.baseURL, c.apiVersion, path)

	if len(query) > 0 {
		values := neturl.Values{}
		for key, value := range query {
			if value != "" {
				values.Set(key, value)
			}
		}
		if encoded := values.Encode(); encoded != "" {
			url += "?" + encoded
		}
	}

	var reader io.Reader
	if body != nil {
		payload, _ := json.Marshal(body)
		reader = bytes.NewReader(payload)
	}

	req, _ := http.NewRequest(method, url, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
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
			return c.do(method, path, body, query, true)
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
