package vault

import (
	"fmt"

	"github.com/ZaiiiRan/messenger/backend/auth-service/internal/config/settings"
	vaultapi "github.com/hashicorp/vault/api"
)

type Client struct {
	client *vaultapi.Client
	path   string
}

func New(cfg settings.VaultSettings) (*Client, error) {
	vaultCfg := vaultapi.DefaultConfig()
	vaultCfg.Address = cfg.Address

	client, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("create vault client: %w", err)
	}

	client.SetToken(cfg.Token)

	health, err := client.Sys().Health()
	if err != nil {
		return nil, fmt.Errorf("ping vault: %w", err)
	}
	if health.Sealed {
		return nil, fmt.Errorf("vault is sealed")
	}

	return &Client{client: client, path: cfg.Path}, nil
}

func (c *Client) GetSecrets() (map[string]string, error) {
	secret, err := c.client.Logical().Read(c.path)
	if err != nil {
		return nil, fmt.Errorf("read secret from vault: %w", err)
	}
	if secret == nil {
		return nil, nil
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		data = secret.Data
	}

	result := make(map[string]string, len(data))
	for k, v := range data {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}

	return result, nil
}
