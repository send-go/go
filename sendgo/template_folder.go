package sendgo

// TemplateFolderService는 기업 계정의 템플릿 공용 폴더를 관리합니다. v2 전용.
type TemplateFolderService struct{ http *httpClient }

// TemplateFolderListQuery는 폴더 트리의 템플릿 수 조회 조건입니다.
type TemplateFolderListQuery struct {
	TemplateType   string
	KakaoSenderKey string
}

// TemplateFolderCreateRequest는 폴더 생성 요청입니다.
type TemplateFolderCreateRequest struct {
	Name       string  `json:"name"`
	ParentUUID *string `json:"parentUuid,omitempty"`
}

// TemplateFolderAssignRequest는 1~100개 템플릿 이동 요청입니다.
type TemplateFolderAssignRequest struct {
	TemplateType   string   `json:"templateType"`
	KakaoSenderKey string   `json:"kakaoSenderKey"`
	TemplateCodes  []string `json:"templateCodes"`
	// nil도 반드시 JSON null로 전송하여 미분류로 이동합니다.
	FolderUUID *string `json:"folderUuid"`
}

// List는 폴더 트리를 조회합니다. TemplateType은 notice 또는 brand입니다.
func (s *TemplateFolderService) List(q TemplateFolderListQuery) (map[string]any, error) {
	query := map[string]string{}
	if q.TemplateType != "" {
		query["templateType"] = q.TemplateType
	}
	if q.KakaoSenderKey != "" {
		query["kakaoSenderKey"] = q.KakaoSenderKey
	}
	return s.http.get("template-folders", query)
}

// Create는 루트 또는 하위 폴더를 생성합니다.
func (s *TemplateFolderService) Create(req TemplateFolderCreateRequest) (map[string]any, error) {
	return s.http.post("template-folders", req)
}

// Assign은 템플릿을 폴더 또는 미분류로 이동합니다.
func (s *TemplateFolderService) Assign(req TemplateFolderAssignRequest) (map[string]any, error) {
	return s.http.patch("template-folders/templates", req)
}
