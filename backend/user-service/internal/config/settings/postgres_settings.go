package settings

import "github.com/spf13/viper"

type PostgresSettings struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	Database string `mapstructure:"database"`
	Options  string `mapstructure:"options"`

	MinConns          uint `mapstructure:"min_conns"`
	MaxConns          uint `mapstructure:"max_conns"`
	MinIdleConns      uint `mapstructure:"min_idle_conns"`
	MaxConnIdleTime   uint `mapstructure:"max_conn_idle_time"`
	MaxConnLifetime   uint `mapstructure:"max_conn_lifetime"`
	PingTimeout       uint `mapstructure:"ping_timeout"`
	HealthCheckPeriod uint `mapstructure:"health_check_period"`
}

func SetPostgresDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".user", "postgres")
	v.SetDefault(prefix+".password", "")
	v.SetDefault(prefix+".address", "localhost:5432")
	v.SetDefault(prefix+".database", "postgres")
	v.SetDefault(prefix+".options", "sslmode=disable")

	v.SetDefault(prefix+".min_conns", 2)
	v.SetDefault(prefix+".max_conns", 20)
	v.SetDefault(prefix+".min_idle_conns", 2)
	v.SetDefault(prefix+".max_conn_idle_time", 900)
	v.SetDefault(prefix+".max_conn_lifetime", 3600)
	v.SetDefault(prefix+".ping_timeout", 5)
	v.SetDefault(prefix+".health_check_period", 30)
}
