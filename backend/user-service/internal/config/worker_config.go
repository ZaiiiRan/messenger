package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	DB                                 settings.PostgresSettings                           `mapstructure:"db"`
	Redis                              settings.RedisSettings                              `mapstructure:"redis"`
	Shutdown                           settings.ShutdownSettings                           `mapstructure:"shutdown"`
	Vault                              settings.VaultSettings                              `mapstructure:"vault"`
	UnconfirmedUsersDataClearingWorker settings.UnconfirmedUsersDataClearingWorkerSettings `mapstructure:"unconfirmed_users_data_clearing_worker"`
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/user-service")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	setWorkerDefaults(v)

	if err := vault.LoadVaultSecrets(v, "vault"); err != nil {
		return nil, err
	}

	var cfg WorkerConfig

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setWorkerDefaults(v *viper.Viper) {
	settings.SetPostgresDefaults(v, "db")
	settings.SetRedisDefaults(v, "redis")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
	settings.SetUnconfirmedUsersDataClearingWorkerDefaults(v, "unconfirmed_users_data_clearing_worker")
}
