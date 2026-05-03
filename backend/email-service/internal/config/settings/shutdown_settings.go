package settings

import "github.com/spf13/viper"

type ShutdownSettings struct {
	ShutdownTimeout uint `mapstructure:"shutdown_timeout"`
}

func SetShutdownDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".shutdown_timeout", 5)
}
