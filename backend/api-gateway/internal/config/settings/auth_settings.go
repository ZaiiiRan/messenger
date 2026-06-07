package settings

import "github.com/spf13/viper"

// AuthSettings configures how the gateway manages the refresh token cookie,
// mirroring the behaviour of the legacy (old) backend.
type AuthSettings struct {
	RefreshTokenCookieName string `mapstructure:"refresh_token_cookie_name"`
	RefreshTokenTTLHours   uint   `mapstructure:"refresh_token_ttl_hours"`
	CookiePath             string `mapstructure:"cookie_path"`
	CookieDomain           string `mapstructure:"cookie_domain"`
	CookieSecure           bool   `mapstructure:"cookie_secure"`
	CookieSameSite         string `mapstructure:"cookie_same_site"`
}

func SetAuthDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".refresh_token_cookie_name", "refreshToken")
	v.SetDefault(prefix+".refresh_token_ttl_hours", 24*30)
	v.SetDefault(prefix+".cookie_path", "/")
	v.SetDefault(prefix+".cookie_domain", "")
	v.SetDefault(prefix+".cookie_secure", false)
	v.SetDefault(prefix+".cookie_same_site", "Lax")
}
