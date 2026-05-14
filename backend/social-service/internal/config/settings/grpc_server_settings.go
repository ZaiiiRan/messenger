package settings

import "github.com/spf13/viper"

type GRPCServerSettings struct {
	Port string `mapstructure:"port"`

	MaxConnectionIdle     uint `mapstructure:"max_connection_idle"`
	MaxConnectionAge      uint `mapstructure:"max_connection_age"`
	MaxConnectionAgeGrace uint `mapstructure:"max_connection_age_grace"`
	KeepaliveTime         uint `mapstructure:"keepalive_time"`
	KeepaliveTimeout      uint `mapstructure:"keepalive_timeout"`

	PermitWithoutStream bool `mapstructure:"permit_without_stream"`
}

func SetGRPCServerDefaults(v *viper.Viper, prefix string, defaultPort string) {
	v.SetDefault(prefix+".port", defaultPort)
	v.SetDefault(prefix+".max_connection_idle", 300)
	v.SetDefault(prefix+".max_connection_age", 0)
	v.SetDefault(prefix+".max_connection_age_grace", 0)
	v.SetDefault(prefix+".keepalive_time", 7200)
	v.SetDefault(prefix+".keepalive_timeout", 20)
	v.SetDefault(prefix+".permit_without_stream", false)
}
