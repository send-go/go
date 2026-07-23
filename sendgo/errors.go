package sendgo

import "fmt"

// SendgoError는 Sendgo API 호출 실패 시 반환되는 에러입니다.
type SendgoError struct {
	StatusCode   int
	ErrorCode    string
	Message      string
	Endpoint     string
	APIVersion   string
	ResponseBody map[string]any
}

func (e *SendgoError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("HTTP %d [%s] %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("HTTP %d %s", e.StatusCode, e.Message)
}

func newSendgoError(status int, body map[string]any, endpoint, apiVersion string) *SendgoError {
	errorCode, _ := body["code"].(string)
	message, _ := body["message"].(string)
	if message == "" {
		message = "Unknown error"
	}
	return &SendgoError{
		StatusCode:   status,
		ErrorCode:    errorCode,
		Message:      message,
		Endpoint:     endpoint,
		APIVersion:   apiVersion,
		ResponseBody: body,
	}
}
