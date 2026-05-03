package settings

import "github.com/spf13/viper"

type VaultSettings struct {
	Address string `mapstructure:"address"`
	Token   string `mapstructure:"token"`
	Path    string `mapstructure:"path"`
	Enabled bool   `mapstructure:"enabled"`
}

func SetVaultDefaults(v *viper.Viper, prefix string) {
	v.SetDefault(prefix+".address", "http://localhost:8200")
	v.SetDefault(prefix+".token", "")
	v.SetDefault(prefix+".path", "secret/data/user-service")
	v.SetDefault(prefix+".enabled", false)
}
