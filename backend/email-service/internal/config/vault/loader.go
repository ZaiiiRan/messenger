package vault

import (
	"strings"

	"github.com/ZaiiiRan/messenger/backend/email-service/internal/config/settings"
	"github.com/ZaiiiRan/messenger/backend/email-service/internal/transport/vault"
	"github.com/spf13/viper"
)

func LoadVaultSecrets(v *viper.Viper, prefix string) error {
	cfg := settings.VaultSettings{
		Address: v.GetString(prefix + ".address"),
		Token:   v.GetString(prefix + ".token"),
		Path:    v.GetString(prefix + ".path"),
		Enabled: v.GetBool(prefix + ".enabled"),
	}

	if !cfg.Enabled {
		return nil
	}

	vc, err := vault.New(cfg)
	if err != nil {
		return err
	}

	secrets, err := vc.GetSecrets()
	if err != nil {
		return err
	}

	injectToViper(v, secrets)
	return nil
}

func injectToViper(v *viper.Viper, secrets map[string]string) {
	replacer := strings.NewReplacer(".", "_")

	for k, val := range secrets {
		normalizedKey := strings.ToLower(replacer.Replace(k))

		if !v.IsSet(normalizedKey) {
			v.SetDefault(normalizedKey, val)
		}
	}
}
