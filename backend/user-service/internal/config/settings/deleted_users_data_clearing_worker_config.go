package settings

import "github.com/spf13/viper"

type DeletedUsersDataClearingWorkerSettings struct {
	Count            uint `mapstructure:"count"`
	IntervalMS       uint `mapstructure:"interval_ms"`
	NoDataIntervalMS uint `mapstructure:"no_data_interval_ms"`
	BatchSize        uint `mapstructure:"batch_size"`
}

func SetDeletedUsersDataClearingWorkerDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".count", 5)
	v.SetDefault(prefix+".interval_ms", 1000)
	v.SetDefault(prefix+".no_data_interval_ms", 900000)
	v.SetDefault(prefix+".batch_size", 1000)
}
