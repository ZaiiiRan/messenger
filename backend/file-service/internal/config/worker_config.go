package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/file-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	DB                         settings.PostgresSettings                   `mapstructure:"db"`
	Redis                      settings.RedisSettings                      `mapstructure:"redis"`
	Shutdown                   settings.ShutdownSettings                   `mapstructure:"shutdown"`
	MetricsServer              settings.MetricsServerSettings              `mapstructure:"metrics_server"`
	Vault                      settings.VaultSettings                      `mapstructure:"vault"`
	UserPhotosS3               settings.S3Settings                         `mapstructure:"user_photos_s3"`
	ChatFilesS3                settings.S3Settings                         `mapstructure:"chat_files_s3"`
	FileDeletionFromS3Worker   settings.FileDeletionFromS3WorkerSettings   `mapstructure:"file_deletion_from_s3_worker"`
	NotInUseFileDeletionWorker settings.NotInUseFileDeletionWorkerSettings `mapstructure:"not_in_use_file_deletion_worker"`
	FileDeletionTasksConsumer  settings.FileDeletionTasksConsumerSettings  `mapstructure:"file_deletion_tasks_consumer"`
	FileDeletionWorker         settings.FileDeletionWorkerSettings         `mapstructure:"file_deletion_worker"`
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/file-service")

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
	settings.SetMetricsServerDefaults(v, "metrics_server")
	settings.SetVaultDefaults(v, "vault")
	settings.SetS3Defaults(v, "user_photos_s3")
	settings.SetS3Defaults(v, "chat_files_s3")
	settings.SetFileDeletionFromS3WorkerDefaults(v, "file_deletion_from_s3_worker")
	settings.SetNotInUseFileDeletionWorkerDefaults(v, "not_in_use_file_deletion_worker")
	settings.SetFileDeletionTasksConsumerSettings(v, "file_deletion_tasks_consumer")
	settings.SetFileDeletionWorkerDefaults(v, "file_deletion_worker")
}
