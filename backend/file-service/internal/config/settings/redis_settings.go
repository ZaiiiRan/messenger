package settings

import "github.com/spf13/viper"

type RedisSettings struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`

	MaxPoolSize     uint `mapstructure:"max_pool_size"`
	MinPoolSize     uint `mapstructure:"min_pool_size"`
	MaxConnIdleTime uint `mapstructure:"max_conn_idle_time"`

	DialTimeout  uint `mapstructure:"dial_timeout"`
	ReadTimeout  uint `mapstructure:"read_timeout"`
	WriteTimeout uint `mapstructure:"write_timeout"`
}

func SetRedisDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".address", "localhost:6379")
	v.SetDefault(prefix+".password", "")
	v.SetDefault(prefix+".min_pool_size", 5)
	v.SetDefault(prefix+".max_pool_size", 200)
	v.SetDefault(prefix+".max_conn_idle_time", 60)
	v.SetDefault(prefix+".dial_timeout", 1000)
	v.SetDefault(prefix+".read_timeout", 200)
	v.SetDefault(prefix+".write_timeout", 200)
}
