package vault

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type Client struct {
	Secrets   vault.Secrets
	MountPath string
}

type Secret struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Cert         string `json:"cert"`
	CertPassword string `json:"cert_password"`
}

func NewClient() *Client {
	c, err := vault.New(
		vault.WithAddress(os.Getenv("VAULT_ADDR")),
		vault.WithRequestTimeout(30*time.Second),
	)
	s, err := c.Auth.AppRoleLogin(
		context.Background(),
		schema.AppRoleLoginRequest{
			RoleId:   os.Getenv("VAULT_ROLE_ID"),
			SecretId: os.Getenv("VAULT_SECRET_ID"),
		},
	)
	err = c.SetToken(s.Auth.ClientToken)
	if err != nil {
		log.Fatal(err)
	}

	vc := &Client{
		Secrets:   c.Secrets,
		MountPath: os.Getenv("VAULT_SECRET_MOUNT_PATH"),
	}

	return vc
}

func (c Client) Get(p string) (*Secret, error) {
	s, err := c.Secrets.KvV2Read(context.Background(), p, vault.WithMountPath(c.MountPath))
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(s.Data.Data)
	if err != nil {
		return nil, err
	}

	var d Secret
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}
