package settings

import "github.com/spf13/viper"

type KafkaSettings struct {
	Brokers  []string `mapstructure:"brokers"`
	User     string   `mapstructure:"user"`
	Password string   `mapstructure:"password"`

	DialTimeout  uint `mapstructure:"dial_timeout"`
}

func SetKafkaDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".brokers", []string{"localhost:9092"})
	v.SetDefault(prefix+".user", "")
	v.SetDefault(prefix+".password", "")
	v.SetDefault(prefix+".dial_timeout", 10)
}
