package sendgo

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
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

func (c *httpClient) put(path string, body any) (map[string]any, error) {
	return c.do(http.MethodPut, path, body, nil, false)
}

func (c *httpClient) patch(path string, body any) (map[string]any, error) {
	return c.do(http.MethodPatch, path, body, nil, false)
}

// delete performs a DELETE request. `do` already drives the verb, so this only
// needs to omit the body.
func (c *httpClient) delete(path string) (map[string]any, error) {
	return c.do(http.MethodDelete, path, nil, nil, false)
}

// MultipartFile은 multipart 업로드에 붙일 파일입니다.
//
// 서류·이미지 첨부가 있는 관리 API(발신번호 등록, 이미지 템플릿, 검수
// 첨부)는 JSON으로 보낼 수 없습니다. 서버 검증이 확장자를 보므로
// FileName은 반드시 채워야 합니다.
type MultipartFile struct {
	// FieldName은 폼 필드 이름입니다 (예: "csuCertificate").
	FieldName string
	// FileName은 서버에 알릴 파일명입니다 (예: "csu.pdf").
	FileName string
	// Content는 파일 내용입니다.
	Content io.Reader
}

// postMultipart는 multipart/form-data POST를 보냅니다.
//
// multipart에는 배열도 불리언도 없으므로, fields의 슬라이스·맵 값은 JSON
// 문자열로 눌러 보냅니다 — 서버가 그렇게 받아 읽습니다.
func (c *httpClient) postMultipart(path string, fields map[string]any, files []MultipartFile) (map[string]any, error) {
	return c.doMultipart(path, fields, files, false)
}

func (c *httpClient) doMultipart(path string, fields map[string]any, files []MultipartFile, isRetry bool) (map[string]any, error) {
	token, err := c.tokenMgr.getToken()
	if err != nil {
		return nil, err
	}

	// 재시도 시 Content를 다시 읽어야 하므로 미리 메모리에 담아 둔다.
	// io.Reader는 한 번 소진되면 되감을 수 없어, 그러지 않으면 토큰 갱신 후
	// 재시도가 빈 파일을 올린다.
	buffered := make([][]byte, len(files))
	for i, file := range files {
		if file.Content == nil {
			continue
		}
		data, readErr := io.ReadAll(file.Content)
		if readErr != nil {
			return nil, fmt.Errorf("sendgo: failed to read %s: %w", file.FieldName, readErr)
		}
		buffered[i] = data
	}

	return c.sendMultipart(path, fields, files, buffered, token, isRetry)
}

func (c *httpClient) sendMultipart(
	path string,
	fields map[string]any,
	files []MultipartFile,
	buffered [][]byte,
	token string,
	isRetry bool,
) (map[string]any, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for key, value := range fields {
		if value == nil {
			continue
		}

		var encoded string
		switch typed := value.(type) {
		case string:
			encoded = typed
		case bool:
			if typed {
				encoded = "1"
			} else {
				encoded = "0"
			}
		case []byte:
			encoded = string(typed)
		default:
			marshalled, marshalErr := json.Marshal(typed)
			if marshalErr != nil {
				return nil, fmt.Errorf("sendgo: failed to encode field %s: %w", key, marshalErr)
			}
			// 숫자·문자열은 JSON 인용부호 없이 보내야 서버가 스칼라로 읽는다.
			encoded = strings.Trim(string(marshalled), `"`)
		}

		if writeErr := writer.WriteField(key, encoded); writeErr != nil {
			return nil, fmt.Errorf("sendgo: failed to write field %s: %w", key, writeErr)
		}
	}

	for i, file := range files {
		part, partErr := writer.CreateFormFile(file.FieldName, file.FileName)
		if partErr != nil {
			return nil, fmt.Errorf("sendgo: failed to create part %s: %w", file.FieldName, partErr)
		}
		if _, writeErr := part.Write(buffered[i]); writeErr != nil {
			return nil, fmt.Errorf("sendgo: failed to write %s: %w", file.FieldName, writeErr)
		}
	}

	if closeErr := writer.Close(); closeErr != nil {
		return nil, fmt.Errorf("sendgo: failed to close multipart writer: %w", closeErr)
	}

	url := fmt.Sprintf("%s/api/%s/%s", c.baseURL, c.apiVersion, path)
	req, _ := http.NewRequest(http.MethodPost, url, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.bearerAuth(token))

	// 파일 업로드는 JSON 요청보다 오래 걸린다.
	uploader := &http.Client{Timeout: 60 * time.Second}

	resp, err := uploader.Do(req)
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
			newToken, tokenErr := c.tokenMgr.getToken()
			if tokenErr != nil {
				return nil, tokenErr
			}
			return c.sendMultipart(path, fields, files, buffered, newToken, true)
		}
		return nil, newSendgoError(resp.StatusCode, responseBody, endpoint, c.apiVersion)
	}

	return responseBody, nil
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
