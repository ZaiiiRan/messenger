package settings

import "github.com/spf13/viper"

type MetricsServerSettings struct {
	Port string `mapstructure:"port"`
}

func SetMetricsServerDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".port", ":8080")
}
