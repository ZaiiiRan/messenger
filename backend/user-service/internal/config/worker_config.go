package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/user-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	MetricsServer                        settings.MetricsServerSettings                        `mapstructure:"metrics_server"`
	DB                                   settings.PostgresSettings                             `mapstructure:"db"`
	Redis                                settings.RedisSettings                                `mapstructure:"redis"`
	FileServiceGRPCClient                settings.GRPCClientSettings                           `mapstructure:"file_service_grpc_client"`
	UserDataDeletionTasksProducer        settings.KafkaProducerSettings                        `mapstructure:"user_data_deletion_tasks_producer"`
	Shutdown                             settings.ShutdownSettings                             `mapstructure:"shutdown"`
	Vault                                settings.VaultSettings                                `mapstructure:"vault"`
	UnconfirmedUsersDataClearingWorker   settings.UnconfirmedUsersDataClearingWorkerSettings   `mapstructure:"unconfirmed_users_data_clearing_worker"`
	DeletedUsersDataClearingWorker       settings.DeletedUsersDataClearingWorkerSettings       `mapstructure:"deleted_users_data_clearing_worker"`
	UnbanTemporarilyBannedUsersWorker    settings.UnbanTemporarilyBannedUsersWorkerSettings    `mapstructure:"unban_temporarily_banned_users_worker"`
	UserDataDeletionTasksSendingWorker   settings.KafkaSendingWorkerSettings                   `mapstructure:"user_data_deletion_tasks_sending_worker"`
	PhotoFileDeletionWorker              settings.KafkaSendingWorkerSettings                   `mapstructure:"photo_file_deletion_worker"`
	UserRelationshipChangesTasksConsumer settings.UserRelationshipChangesTasksConsumerSettings `mapstructure:"user_relationship_changes_tasks_consumer"`
	UpdateFriendsInPrivacySettingsWorker settings.UpdateFriendsInPrivacySettingsWorker         `mapstructure:"update_friends_in_privacy_settings_worker"`
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
	settings.SetMetricsServerDefaults(v, "metrics_server")
	settings.SetPostgresDefaults(v, "db")
	settings.SetRedisDefaults(v, "redis")
	settings.SetGRPCClientDefaults(v, "file_service_grpc_client", "localhost:50054")
	settings.SetKafkaProducerDefaults(v, "user_data_deletion_tasks_producer", "user-data-deletion-tasks")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
	settings.SetUnconfirmedUsersDataClearingWorkerDefaults(v, "unconfirmed_users_data_clearing_worker")
	settings.SetDeletedUsersDataClearingWorkerDefaults(v, "deleted_users_data_clearing_worker")
	settings.SetUnbanTemporarilyBannedUsersWorkerDefaults(v, "unban_temporarily_banned_users_worker")
	settings.SetKafkaSendingWorkerDefaults(v, "user_data_deletion_tasks_sending_worker")
	settings.SetKafkaSendingWorkerDefaults(v, "photo_file_deletion_worker")
	settings.SetUserRelationshipChangesTasksConsumerSettings(v, "user_relationship_changes_tasks_consumer")
	settings.SetUpdateFriendsInPrivacySettingsWorkerDefaults(v, "update_friends_in_privacy_settings_worker")
}
