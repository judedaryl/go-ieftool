package msgraph

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type Client struct {
	GraphServiceClient     *sdk.GraphServiceClient
	ClientSecretCredential *azidentity.ClientSecretCredential
	Scopes                 []string
	Token                  azcore.AccessToken
}

func NewClient(tid, cid, s string) (*Client, error) {
	g := &Client{
		Scopes: []string{"https://graph.microsoft.com/.default"},
	}

	cr, err := azidentity.NewClientSecretCredential(
		tid,
		cid,
		s,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create client credentials: %s", err.Error())
	}
	g.ClientSecretCredential = cr
	t, err := g.ClientSecretCredential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: g.Scopes,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get token: %s", err.Error())
	}
	g.Token = t

	c, err := sdk.NewGraphServiceClientWithCredentials(cr, g.Scopes)
	if err != nil {
		return nil, err
	}
	g.GraphServiceClient = c

	return g, nil
}
