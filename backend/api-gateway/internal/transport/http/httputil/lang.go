package httputil

import "strings"

// extractLang normalizes an Accept-Language header value down to a bare
// language tag (e.g. "ru-RU,en;q=0.9" -> "ru").
func extractLang(acceptLang string) string {
	s := strings.SplitN(acceptLang, ",", 2)[0]
	s = strings.TrimSpace(strings.SplitN(s, ";", 2)[0])
	if idx := strings.IndexByte(s, '-'); idx != -1 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}
