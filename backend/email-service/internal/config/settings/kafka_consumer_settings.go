package settings

import "github.com/spf13/viper"

type KafkaConsumerSettings struct {
	KafkaSettings  KafkaSettings `mapstructure:"kafka_settings"`
	GroupID        string        `mapstructure:"group_id"`
	Topic          string        `mapstructure:"topic"`
	BatchSize      uint          `mapstructure:"batch_size"`
	BatchTimeoutMs uint          `mapstructure:"batch_timeout_ms"`
}

func SetKafkaConsumerDefaults(v *viper.Viper, prefix string, defaultTopic string) {
	SetKafkaDefaults(v, prefix+".kafka_settings")
	v.SetDefault(prefix+".group_id", "email-service-consumer")
	v.SetDefault(prefix+".topic", defaultTopic)
	v.SetDefault(prefix+".batch_size", 100)
	v.SetDefault(prefix+".batch_timeout_ms", 1000)
}
