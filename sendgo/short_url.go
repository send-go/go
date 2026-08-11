package sendgo

import "net/url"

// ShortURLService는 짧은 URL 서비스입니다.
//
// 메시지에 넣는 링크를 줄이고 클릭 반응을 집계합니다. v2 전용입니다.
//
// 사용법:
//
//	created, err := client.ShortURL.Create(sendgo.ShortURLRequest{
//	    TargetURL: "https://example.com/promotions/summer-sale",
//	    Title:     "여름 세일 랜딩",
//	})
//
//	// created["data"]["shortUrl"] 를 문자/알림톡 본문에 넣습니다.
type ShortURLService struct {
	http *httpClient
}

func newShortURLService(hc *httpClient) *ShortURLService {
	return &ShortURLService{http: hc}
}

// Create는 짧은 URL을 만듭니다.
//
// 같은 원본 URL을 다시 줄이면 기존 링크가 그대로 반환됩니다.
// 캠페인별로 반응을 분리해 집계하려면 ForceNew를 true로 둡니다.
func (s *ShortURLService) Create(req ShortURLRequest) (map[string]any, error) {
	return s.http.post("short-urls", req)
}

// List는 짧은 URL 목록을 조회합니다.
func (s *ShortURLService) List(query ShortURLListQuery) (map[string]any, error) {
	return s.http.get("short-urls", query.toMap())
}

// Show는 짧은 URL 상세를 조회합니다.
func (s *ShortURLService) Show(code string) (map[string]any, error) {
	return s.http.get("short-urls/"+url.PathEscape(code), nil)
}

// Stats는 반응 통계를 조회합니다.
// 일별 추이와 디바이스/유입경로/국가별 분해가 담깁니다.
func (s *ShortURLService) Stats(code string, query ShortURLStatsQuery) (map[string]any, error) {
	return s.http.get("short-urls/"+url.PathEscape(code)+"/stats", query.toMap())
}

// Deactivate는 리다이렉트를 중지합니다.
//
// 링크는 삭제되지 않고 누적 통계도 남습니다. 이후 그 링크로 들어오면
// 410 Gone 이 반환됩니다.
func (s *ShortURLService) Deactivate(code string) (map[string]any, error) {
	return s.http.delete("short-urls/" + url.PathEscape(code))
}
