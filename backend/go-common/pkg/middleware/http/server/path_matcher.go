package middleware

type PathMatcher func(method, path string) bool

func ProtectAll() PathMatcher {
	return func(_, _ string) bool { return true }
}
