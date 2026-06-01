package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	MetricsServer         settings.MetricsServerSettings `mapstructure:"metrics_server"`
	FileServiceGRPCClient settings.GRPCClientSettings    `mapstructure:"file_service_grpc_client"`
	DB                    settings.PostgresSettings      `mapstructure:"db"`
	Shutdown              settings.ShutdownSettings      `mapstructure:"shutdown"`
	Vault                 settings.VaultSettings         `mapstructure:"vault"`
	UserPhotosS3          settings.S3Settings            `mapstructure:"user_photos_s3"`
	ChatFilesS3           settings.S3Settings            `mapstructure:"chat_files_s3"`

	Processor settings.ProcessorSettings `mapstructure:"processor"`

	FileProcessingWorker settings.FileProcessingWorkerSettings `mapstructure:"file_processing_worker"`

	FileProcessingCompletedEventsProducer settings.KafkaProducerSettings      `mapstructure:"file_processing_completed_events_producer"`
	FileProcessingCompletedSendingWorker  settings.KafkaSendingWorkerSettings `mapstructure:"file_processing_completed_sending_worker"`
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/media-service")

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
	settings.SetGRPCClientDefaults(v, "file_service_grpc_client", "localhost:50052")
	settings.SetPostgresDefaults(v, "db")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
	settings.SetS3Defaults(v, "user_photos_s3")
	settings.SetS3Defaults(v, "chat_files_s3")
	settings.SetProcessorDefaults(v, "processor")
	settings.SetFileProcessingWorkerDefaults(v, "file_processing_worker")
	settings.SetKafkaProducerDefaults(v, "file_processing_completed_events_producer", "media-processing-completed")
	settings.SetKafkaSendingWorkerDefaults(v, "file_processing_completed_sending_worker")
}
