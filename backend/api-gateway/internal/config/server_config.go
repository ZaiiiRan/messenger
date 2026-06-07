package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/api-gateway/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	HTTPServer    settings.HTTPServerSettings    `mapstructure:"http_server"`
	MetricsServer settings.MetricsServerSettings `mapstructure:"metrics_server"`

	AuthServiceGRPCClient   settings.GRPCClientSettings `mapstructure:"auth_service_grpc_client"`
	UserServiceGRPCClient   settings.GRPCClientSettings `mapstructure:"user_service_grpc_client"`
	SocialServiceGRPCClient settings.GRPCClientSettings `mapstructure:"social_service_grpc_client"`
	FileServiceGRPCClient   settings.GRPCClientSettings `mapstructure:"file_service_grpc_client"`
	MediaServiceGRPCClient  settings.GRPCClientSettings `mapstructure:"media_service_grpc_client"`

	JWT      settings.JWTSettings      `mapstructure:"jwt"`
	Auth     settings.AuthSettings     `mapstructure:"auth"`
	CORS     settings.CORSSettings     `mapstructure:"cors"`
	Swagger  settings.SwaggerSettings  `mapstructure:"swagger"`
	Shutdown settings.ShutdownSettings `mapstructure:"shutdown"`
	Vault    settings.VaultSettings    `mapstructure:"vault"`
}

func LoadServerConfig() (*ServerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/api-gateway")

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
	settings.SetHTTPServerDefaults(v, "http_server", ":8000")
	settings.SetMetricsServerDefaults(v, "metrics_server")

	settings.SetGRPCClientDefaults(v, "auth_service_grpc_client", "localhost:50052")
	settings.SetGRPCClientDefaults(v, "user_service_grpc_client", "localhost:50051")
	settings.SetGRPCClientDefaults(v, "social_service_grpc_client", "localhost:50053")
	settings.SetGRPCClientDefaults(v, "file_service_grpc_client", "localhost:50054")
	settings.SetGRPCClientDefaults(v, "media_service_grpc_client", "localhost:50055")

	settings.SetJWTDefaults(v, "jwt")
	settings.SetAuthDefaults(v, "auth")
	settings.SetCORSDefaults(v, "cors")
	settings.SetSwaggerDefaults(v, "swagger")
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
}
