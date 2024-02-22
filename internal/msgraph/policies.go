package msgraph

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"com.schumann-it.go-ieftool/internal/msgraph/trustframework"
)

func (c *Client) ListPolicies() ([]string, error) {
	r, err := c.GraphServiceClient.TrustFramework().Policies().Get(context.Background(), nil)
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

func (c *Client) DeletePolicies() error {
	ps, err := c.ListPolicies()
	if err != nil {
		return err
	}

	for _, id := range ps {
		err = c.GraphServiceClient.TrustFramework().Policies().ByTrustFrameworkPolicyId(id).Delete(context.Background(), nil)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to delete trustframework.Policy %s: %s", id, err))
			continue
		}
		log.Println(fmt.Sprintf("Policy %s deleted", id))
	}

	return nil
}

func (c *Client) UploadPolicies(policies []trustframework.Policy) {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	for _, p := range policies {
		go c.uploadPolicy(p, &wg)
	}
	wg.Wait()
}

func (c *Client) uploadPolicy(p trustframework.Policy, wg *sync.WaitGroup) {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.Token))
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
		log.Fatalf("Upload failed for trustframework.Policy %s \n%s\n", p.PolicyId, string(body))
	}

	log.Println(fmt.Sprintf("Policy %s uploaded", p.PolicyId))
}
