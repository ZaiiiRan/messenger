package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/social-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/social-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	MetricsServer                             settings.MetricsServerSettings      `mapstructure:"metrics_server"`
	DB                                        settings.PostgresSettings           `mapstructure:"db"`
	Redis                                     settings.RedisSettings              `mapstructure:"redis"`
	UserRelationshipChangesTasksProducer      settings.KafkaProducerSettings      `mapstructure:"user_relationship_changes_tasks_producer"`
	Shutdown                                  settings.ShutdownSettings           `mapstructure:"shutdown"`
	Vault                                     settings.VaultSettings              `mapstructure:"vault"`
	UserRelationshipChangesTasksSendindWorker settings.KafkaSendingWorkerSettings `mapstructure:"user_relationship_changes_tasks_sending_worker"`
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/social-service")

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
	settings.SetKafkaProducerDefaults(v, "user_relationship_changes_tasks_producer", "")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
	settings.SetKafkaSendingWorkerDefaults(v, "user_relationship_changes_tasks_sending_worker")
}
