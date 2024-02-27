package msgraph

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
)

func (c *Client) FixAppRegistration(id string, patch []byte) error {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	ep := fmt.Sprintf("https://graph.microsoft.com/beta/applications/%s", id)
	req, err := http.NewRequest(http.MethodPatch, ep, bytes.NewBuffer(patch))
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
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
		log.Fatalf("Patch app failed \n%s\n", string(body))
	}

	return nil
}
