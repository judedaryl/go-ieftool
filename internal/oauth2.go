package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TokenResponse struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

func GetToken(clientId, clientSecret, tenantId string) TokenResponse {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	scope := "https://graph.microsoft.com/.default"
	body := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&scope=%s", clientId, clientSecret, scope)
	endpoint := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantId)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer([]byte(body)))
	Check(err)
	resp, err := client.Do(req)
	Check(err)
	respB, err := ioutil.ReadAll(resp.Body)
	Check(err)
	token := &TokenResponse{}
	json.Unmarshal(respB, token)
	return *token
}
