package settings

import "github.com/spf13/viper"

type JWTSettings struct {
	AccessTokenSecret string `mapstructure:"access_token_secret"`
}

func SetJWTDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".access_token_secret", "access-secret-key")
}
