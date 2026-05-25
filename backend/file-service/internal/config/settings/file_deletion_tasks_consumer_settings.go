package settings

import "github.com/spf13/viper"

type FileDeletionTasksConsumerSettings struct {
	KafkaConsumerSettings KafkaConsumerSettings `mapstructure:"kafka_consumer"`
	Count                 uint                  `mapstructure:"count"`
}

func SetFileDeletionTasksConsumerSettings(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".count", 5)
	SetKafkaConsumerDefaults(v, prefix+".kafka_consumer", "file-deletion-tasks")
}
