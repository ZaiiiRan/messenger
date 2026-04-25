package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	GRPCServer            settings.GRPCServerSettings `mapstructure:"grpc_server"`
	UserServiceGRPCClient settings.GRPCClientSettings `mapstructure:"user_service_grpc_client"`
	DB                    settings.PostgresSettings   `mapstructure:"db"`
	Migrate               settings.MigrateSettings    `mapstructure:"migrate"`
	Redis                 settings.RedisSettings      `mapstructure:"redis"`
	Shutdown              settings.ShutdownSettings   `mapstructure:"shutdown"`
	Vault                 settings.VaultSettings      `mapstructure:"vault"`
}

func LoadServerConfig() (*ServerConfig, error) {
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
	settings.SetGRPCClientDefaults(v, "user_service_grpc_client", "localhost:50052")
	settings.SetPostgresDefaults(v, "db")
	settings.SetMigrateDefaults(v, "migrate")
	settings.SetRedisDefaults(v, "redis")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
}
