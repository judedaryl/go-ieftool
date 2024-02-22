package internal

import (
	"fmt"
	"os"
	"strings"

	"com.schumann-it.go-ieftool/internal/msgraph"
)

func NewGraphClientFromEnvironment(e Environment) (*msgraph.Client, error) {
	es := strings.ReplaceAll(fmt.Sprintf("B2C_CLIENT_SECRET_%s", strings.ToUpper(e.Name)), "-", "_")
	c, err := msgraph.NewClient(e.TenantId, e.ClientId, os.Getenv(es))
	if err != nil {
		return nil, fmt.Errorf("could not create client credentials. Did you send the env var %s?: %s", es, err.Error())
	}

	return c, nil
}
