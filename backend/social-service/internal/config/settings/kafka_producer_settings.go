package settings

import "github.com/spf13/viper"

type KafkaProducerSettings struct {
	KafkaSettings KafkaSettings `mapstructure:"kafka_settings"`
	Topic         string        `mapstructure:"topic"`
	ClientID      string        `mapstructure:"client_id"`
	WriteTimeout  uint          `mapstructure:"write_timeout"`
	Name          string        `mapstructure:"name"`
}

func SetKafkaProducerDefaults(v *viper.Viper, prefix string, defaultTopic string) {
	SetKafkaDefaults(v, prefix+".kafka_settings")
	v.SetDefault(prefix+".topic", defaultTopic)
	v.SetDefault(prefix+".write_timeout", 30)
	v.SetDefault(prefix+".name", defaultTopic+"_producer")
}
