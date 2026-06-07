package settings

import "github.com/spf13/viper"

type SwaggerSettings struct {
	Enabled  bool   `mapstructure:"enabled"`
	BasePath string `mapstructure:"base_path"`
}

func SetSwaggerDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".enabled", true)
	v.SetDefault(prefix+".base_path", "/swagger")
}
