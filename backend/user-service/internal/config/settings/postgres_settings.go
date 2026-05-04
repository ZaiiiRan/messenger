package settings

import "github.com/spf13/viper"

type PostgresSettings struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Address  string `mapstructure:"address"`
	Database string `mapstructure:"database"`
	Options  string `mapstructure:"options"`
}

func SetPostgresDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".user", "postgres")
	v.SetDefault(prefix+".password", "postgres")
	v.SetDefault(prefix+".address", "localhost:5432")
	v.SetDefault(prefix+".database", "postgres")
	v.SetDefault(prefix+".options", "sslmode=disable")
}
