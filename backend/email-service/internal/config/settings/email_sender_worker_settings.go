package settings

import "github.com/spf13/viper"

type EmailSenderWorkerSettings struct {
	KafkaConsumerSettings KafkaConsumerSettings `mapstructure:"kafka_consumer"`
	Count                 int                   `mapstructure:"count"`
}

func SetEmailSenderWorkerDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".count", 5)
	SetKafkaConsumerDefaults(v, prefix+".kafka_consumer", "email-codes-tasks")
}
