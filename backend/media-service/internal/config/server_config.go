package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/media-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	GRPCServer            settings.GRPCServerSettings    `mapstructure:"grpc_server"`
	FileServiceGRPCClient settings.GRPCClientSettings    `mapstructure:"file_service_grpc_client"`
	JWT                   settings.JWTSettings           `mapstructure:"jwt"`
	DB                    settings.PostgresSettings      `mapstructure:"db"`
	Migrate               settings.MigrateSettings       `mapstructure:"migrate"`
	Shutdown              settings.ShutdownSettings      `mapstructure:"shutdown"`
	MetricsServer         settings.MetricsServerSettings `mapstructure:"metrics_server"`
	Vault                 settings.VaultSettings         `mapstructure:"vault"`
	UserPhotosS3          settings.S3Settings            `mapstructure:"user_photos_s3"`
	ChatFilesS3           settings.S3Settings            `mapstructure:"chat_files_s3"`
}

func LoadServerConfig() (*ServerConfig, error) {
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

	setServerDefaults(v)

	if err := vault.LoadVaultSecrets(v, "vault"); err != nil {
		return nil, err
	}

	var cfg ServerConfig

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setServerDefaults(v *viper.Viper) {
	settings.SetGRPCServerDefaults(v, "grpc_server", ":50051")
	settings.SetGRPCClientDefaults(v, "file_service_grpc_client", "localhost:50052")
	settings.SetJWTDefaults(v, "jwt")
	settings.SetPostgresDefaults(v, "db")
	settings.SetMigrateDefaults(v, "migrate")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetMetricsServerDefaults(v, "metrics_server")
	settings.SetVaultDefaults(v, "vault")
	settings.SetS3Defaults(v, "user_photos_s3")
	settings.SetS3Defaults(v, "chat_files_s3")
}
