package middleware

type MethodMatcher func(method string) bool

func MiddlewareOnly(methods ...string) MethodMatcher {
	m := make(map[string]struct{}, len(methods))
	for _, s := range methods {
		m[s] = struct{}{}
	}
	return func(method string) bool {
		_, ok := m[method]
		return ok
	}
}
