package config

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/vault"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type WorkerConfig struct {
	Shutdown          settings.ShutdownSettings          `mapstructure:"shutdown"`
	Vault             settings.VaultSettings             `mapstructure:"vault"`
	SMTPClient        settings.SMTPClientSettings        `mapstructure:"smtp_client"`
	EmailSenderWorker settings.EmailSenderWorkerSettings `mapstructure:"email_sender_worker"`
}

func LoadWorkerConfig() (*WorkerConfig, error) {
	_ = godotenv.Load()

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/email-service")

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
	settings.SetShutdownDefaults(v, "shutdown")
	settings.SetVaultDefaults(v, "vault")
	settings.SetEmailSenderWorkerDefaults(v, "email_sender_worker")
	settings.SetSMTPClientDefaults(v, "smtp_client")
}
