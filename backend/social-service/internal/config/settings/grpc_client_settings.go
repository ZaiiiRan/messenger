package settings

import "github.com/spf13/viper"

type GRPCClientSettings struct {
	Address              string `mapstructure:"address"`
	AutoConnect          bool   `mapstructure:"auto_connect"`
	RetriesCount         uint   `mapstructure:"retries_count"`
	PerCallTimeout       uint   `mapstructure:"per_call_timeout"`
	WaitGRPCReadyTimeout uint   `mapstructure:"wait_grpc_ready_timeout"`

	LBPolicy          string  `mapstructure:"lb_policy"`
	MinConnectTimeout uint    `mapstructure:"min_connect_timeout"`
	BackoffBaseDelay  uint    `mapstructure:"backoff_base_delay"`
	BackoffMultiplier float64 `mapstructure:"backoff_multiplier"`
	BackoffMaxDelay   uint    `mapstructure:"backoff_max_delay"`

	KeepaliveTime                uint `mapstructure:"keepalive_time"`
	KeepaliveTimeout             uint `mapstructure:"keepalive_timeout"`
	KeepalivePermitWithoutStream bool `mapstructure:"keepalive_permit_without_stream"`
}

func SetGRPCClientDefaults(v *viper.Viper, prefix string, defaultAdress string) {
	v.SetDefault(prefix+".address", defaultAdress)
	v.SetDefault(prefix+".auto_connect", true)
	v.SetDefault(prefix+".retries_count", 3)
	v.SetDefault(prefix+".per_call_timeout", 2)
	v.SetDefault(prefix+".wait_grpc_ready_timeout", 5)
	v.SetDefault(prefix+".lb_policy", "round_robin")
	v.SetDefault(prefix+".min_connect_timeout", 2)
	v.SetDefault(prefix+".backoff_base_delay", 100)
	v.SetDefault(prefix+".backoff_multiplier", 1.6)
	v.SetDefault(prefix+".backoff_max_delay", 2000)
	v.SetDefault(prefix+".keepalive_time", 0)
	v.SetDefault(prefix+".keepalive_timeout", 0)
	v.SetDefault(prefix+".keepalive_permit_without_stream", false)
}
