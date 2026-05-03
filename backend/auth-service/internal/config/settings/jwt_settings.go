package settings

import "github.com/spf13/viper"

type JWTSettings struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`

	AccessTokenTTL  uint `mapstructure:"access_token_ttl"`
	RefreshTokenTTL uint `mapstructure:"refresh_token_ttl"`
}

func SetJWTDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".access_token_secret", "access-secret-key")
	v.SetDefault(prefix+".refresh_token_secret", "refresh-secret-key")
	v.SetDefault(prefix+".access_token_ttl", 900)    // 15 minutes in seconds
	v.SetDefault(prefix+".refresh_token_ttl", 86400) // 24 hours in seconds (1 day)
}
