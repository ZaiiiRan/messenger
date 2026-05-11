package settings

import "github.com/spf13/viper"

type UserDataDeletionTasksWorkerSettings struct {
	Count           uint `mapstructure:"count"`
	IntervalMS      uint `mapstructure:"interval_ms"`
	RetryIntervalMS uint `mapstructure:"retry_interval_ms"`
	BatchSize       uint `mapstructure:"batch_size"`
}

func SetUserDataDeletionTasksWorkerDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".count", 5)
	v.SetDefault(prefix+".interval_ms", 1000)
	v.SetDefault(prefix+".retry_interval_ms", 300000)
	v.SetDefault(prefix+".batch_size", 100)
}
