package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func UploadPolicy(accessToken, policyId string, content []byte) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	endpoint := fmt.Sprintf("https://graph.microsoft.com/beta/trustFramework/policies/%s/$value", policyId)
	req, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(content))
	if err != nil {
		panic(err)
	}

	authToken := fmt.Sprintf("Bearer %s", accessToken)
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")
	req.Header.Set("Authorization", authToken)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode >= 400 {
		log.Fatalf("Upload failed\n%s\n", string(body))
	}
}
