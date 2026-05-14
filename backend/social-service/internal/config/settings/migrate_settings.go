package settings

import "github.com/spf13/viper"

type MigrateSettings struct {
	NeedToMigrate bool `mapstructure:"need_to_migrate"`
}

func SetMigrateDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".need_to_migrate", false)
}
