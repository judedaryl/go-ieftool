package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type PolicyResponse struct {
	Id string `json:"id"`
}

type PolicyListResponse struct {
	Value []PolicyResponse `json:"value"`
}

type Api struct {
	Token TokenResponse
}

type GraphClient struct {
	c  *msgraphsdk.GraphServiceClient
	cr *azidentity.ClientSecretCredential
	s  []string
	t  azcore.AccessToken
}

func MustNewGraphClientFromEnvironment(e Environment) *GraphClient {
	g, err := NewGraphClientFromEnvironment(e)
	if err != nil {
		log.Fatal(err.Error())
	}

	return g
}
func NewGraphClientFromEnvironment(e Environment) (*GraphClient, error) {
	g := &GraphClient{
		s: []string{"https://graph.microsoft.com/.default"},
	}

	es := strings.ReplaceAll(fmt.Sprintf("B2C_CLIENT_SECRET_%s", strings.ToUpper(e.Name)), "-", "_")
	cr, err := azidentity.NewClientSecretCredential(
		e.TenantId,
		e.ClientId,
		os.Getenv(es),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create client credentials. Did you send the env var %s?: %s", es, err.Error())
	}
	g.cr = cr

	t, err := g.cr.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: g.s,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get token. Did you send the env var %s?: %s", es, err)
	}
	g.t = t

	c, err := msgraphsdk.NewGraphServiceClientWithCredentials(cr, g.s)
	if err != nil {
		return nil, err
	}
	g.c = c

	return g, nil
}

func (g *GraphClient) ListPolicies() ([]string, error) {
	r, err := g.c.TrustFramework().Policies().Get(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	var i []string
	for _, p := range r.GetValue() {
		id := p.GetId()
		i = append(i, *id)
	}

	return i, nil
}

func (g *GraphClient) DeletePolicies() error {
	ps, err := g.ListPolicies()
	if err != nil {
		return err
	}

	for _, id := range ps {
		err = g.c.TrustFramework().Policies().ByTrustFrameworkPolicyId(id).Delete(context.Background(), nil)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to delete policy %s: %s", id, err))
			continue
		}
		log.Println(fmt.Sprintf("Policy %s deleted", id))
	}

	return nil
}

func (g *GraphClient) UploadPolicies(policies []Policy) {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	for _, p := range policies {
		go g.uploadPolicy(p, &wg)
	}
	wg.Wait()
}

func (g *GraphClient) uploadPolicy(p Policy, wg *sync.WaitGroup) {
	defer wg.Done()

	content, _ := os.ReadFile(p.Path)
	client := &http.Client{}
	defer client.CloseIdleConnections()

	ep := fmt.Sprintf("https://graph.microsoft.com/beta/trustFramework/policies/%s/$value", p.PolicyId)
	req, err := http.NewRequest(http.MethodPut, ep, bytes.NewBuffer(content))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", g.t.Token))
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		log.Fatalf("Upload failed for policy %s \n%s\n", p.PolicyId, string(body))
	}

	log.Println(fmt.Sprintf("Policy %s uploaded", p.PolicyId))
}
