package settings

import "github.com/spf13/viper"

type SMTPClientSettings struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaxRetries      uint   `mapstructure:"max_retries"`
	RetryDelayMS    uint   `mapstructure:"retry_delay_ms"`
	RetryMaxDelayMS uint   `mapstructure:"retry_max_delay_ms"`
}

func SetSMTPClientDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".max_retries", 3)
	v.SetDefault(prefix+".retry_delay_ms", 1000)
	v.SetDefault(prefix+".retry_max_delay_ms", 5000)
}
