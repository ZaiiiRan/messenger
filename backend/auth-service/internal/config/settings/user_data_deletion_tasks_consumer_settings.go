package settings

import "github.com/spf13/viper"

type UserDataDeletionTasksConsumerSettings struct {
	KafkaConsumerSettings KafkaConsumerSettings `mapstructure:"kafka_consumer"`
	Count                 uint                  `mapstructure:"count"`
}

func SetUserDataDeletionTasksConsumerDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".count", 5)
	SetKafkaConsumerDefaults(v, prefix+".kafka_consumer", "user-data-deletion-tasks")
}
