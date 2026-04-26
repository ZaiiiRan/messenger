package settings

import "github.com/spf13/viper"

type KafkaProducerSettings struct {
	KafkaSettings  KafkaSettings `mapstructure:"kafka_settings"`
	Topic          string        `mapstructure:"topic"`
	ClientID       string        `mapstructure:"client_id"`
	WriteTimeout   uint          `mapstructure:"write_timeout"`
	BatchSize      uint          `mapstructure:"batch_size"`
	FlushFrequency uint          `mapstructure:"flush_frequency"`
}

func SetKafkaProducerDefaults(v *viper.Viper, prefix string, defaultTopic string) {
	SetKafkaDefaults(v, prefix+".kafka_settings")
	v.SetDefault(prefix+".topic", defaultTopic)
	v.SetDefault(prefix+".write_timeout", 300)
	v.SetDefault(prefix+".batch_size", 100)
	v.SetDefault(prefix+".flush_frequency", 1000)
}
