package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	MetricsServer                           settings.MetricsServerSettings                 `mapstructure:"metrics_server"`
	DB                                      settings.PostgresSettings                      `mapstructure:"db"`
	Redis                                   settings.RedisSettings                         `mapstructure:"redis"`
	Shutdown                                settings.ShutdownSettings                      `mapstructure:"shutdown"`
	Vault                                   settings.VaultSettings                         `mapstructure:"vault"`
	ExpiredTokenClearingWorker              settings.ExpiredTokenClearingWorkerSettings    `mapstructure:"expired_token_clearing_worker"`
	ExpiredResetPasswordCodesClearingWorker settings.ExpiredCodesClearingWorkerSettings    `mapstructure:"expired_reset_password_codes_clearing_worker"`
	ExpiredActivationCodesClearingWorker    settings.ExpiredCodesClearingWorkerSettings    `mapstructure:"expired_activation_codes_clearing_worker"`
	ExpiredEmailChangeCodesClearingWorker   settings.ExpiredCodesClearingWorkerSettings    `mapstructure:"expired_email_change_codes_clearing_worker"`
	UserDataDeletionTasksConsumer           settings.UserDataDeletionTasksConsumerSettings `mapstructure:"user_data_deletion_tasks_consumer"`
	UserDataDeletionTasksWorker             settings.UserDataDeletionTasksWorkerSettings   `mapstructure:"user_data_deletion_tasks_worker"`
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/auth-service")

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
	settings.SetMetricsServerDefaults(v, "metrics_server")
	settings.SetPostgresDefaults(v, "db")
	settings.SetRedisDefaults(v, "redis")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
	settings.SetExpiredTokenClearingWorkerDefaults(v, "expired_token_clearing_worker")
	settings.SetExpiredCodesClearingWorkerDefaults(v, "expired_reset_password_codes_clearing_worker")
	settings.SetExpiredCodesClearingWorkerDefaults(v, "expired_activation_codes_clearing_worker")
	settings.SetExpiredCodesClearingWorkerDefaults(v, "expired_email_change_codes_clearing_worker")
	settings.SetUserDataDeletionTasksConsumerDefaults(v, "user_data_deletion_tasks_consumer")
	settings.SetUserDataDeletionTasksWorkerDefaults(v, "user_data_deletion_tasks_worker")
}
