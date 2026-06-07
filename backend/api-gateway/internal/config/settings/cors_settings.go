package settings

import "github.com/spf13/viper"

type CORSSettings struct {
	Enabled          bool   `mapstructure:"enabled"`
	AllowOrigins     string `mapstructure:"allow_origins"`
	AllowMethods     string `mapstructure:"allow_methods"`
	AllowHeaders     string `mapstructure:"allow_headers"`
	AllowCredentials bool   `mapstructure:"allow_credentials"`
}

func SetCORSDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".enabled", true)
	v.SetDefault(prefix+".allow_origins", "http://localhost:5173")
	v.SetDefault(prefix+".allow_methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	v.SetDefault(prefix+".allow_headers", "Origin,Content-Type,Accept,Authorization,Accept-Language,X-Request-ID")
	v.SetDefault(prefix+".allow_credentials", true)
}
