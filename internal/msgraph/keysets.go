package msgraph

import (
	"context"
	"log"

	"com.schumann-it.go-ieftool/internal/msgraph/trustframework"
	"com.schumann-it.go-ieftool/internal/vault"
	sdkmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	sdktrustframework "github.com/microsoftgraph/msgraph-beta-sdk-go/trustframework"
)

func (c *Client) CreateKeySets(s *vault.Secret) error {
	n := trustframework.NewKeySet([]string{"B2C_1A_TokenSigningKeyContainer", "B2C_1A_TokenEncryptionKeyContainer"})
	if s != nil {
		n.Add("B2C_1A_SamlIdpCert")
	}

	resp, err := c.GraphServiceClient.TrustFramework().KeySets().Get(context.Background(), nil)
	if err != nil {
		log.Fatalln(err)
	}
	for _, ks := range resp.GetValue() {
		id := ks.GetId()
		n.Remove(*id)
	}

	for _, id := range n.IDs {
		switch id {
		case "B2C_1A_TokenSigningKeyContainer":
			err = c.createKeySet(id, "sig")
			if err != nil {
				log.Fatalln(err)
			}
			break
		case "B2C_1A_TokenEncryptionKeyContainer":
			err = c.createKeySet(id, "enc")
			if err != nil {
				log.Fatalln(err)
			}
			break
		case "B2C_1A_SamlIdpCert":
			err = c.uploadCertificate(id, s.Cert, s.CertPassword)
			if err != nil {
				log.Fatalln(err)
			}
			break
		default:
			log.Fatalf("Key Set %s not recognized", id)
		}
	}

	return nil
}

func (c *Client) uploadCertificate(id, cert, pw string) error {
	ks := sdkmodels.NewTrustFrameworkKeySet()
	ks.SetId(&id)
	_, err := c.GraphServiceClient.TrustFramework().KeySets().Post(context.Background(), ks, nil)
	if err != nil {
		return err
	}
	k := sdktrustframework.NewKeySetsItemUploadPkcs12PostRequestBody()
	k.SetKey(&cert)
	k.SetPassword(&pw)
	_, err = c.GraphServiceClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).UploadPkcs12().Post(context.Background(), k, nil)

	return err
}

func (c *Client) createKeySet(id, use string) error {
	ks := sdkmodels.NewTrustFrameworkKeySet()
	ks.SetId(&id)
	_, err := c.GraphServiceClient.TrustFramework().KeySets().Post(context.Background(), ks, nil)
	if err != nil {
		return err
	}
	k := sdktrustframework.NewKeySetsItemGenerateKeyPostRequestBody()
	k.SetUse(&use)
	kty := "RSA"
	k.SetKty(&kty)
	_, err = c.GraphServiceClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(id).GenerateKey().Post(context.Background(), k, nil)

	return err
}

func (c *Client) DeleteKeySets() error {
	a := trustframework.NewKeySet([]string{"B2C_1A_TokenSigningKeyContainer", "B2C_1A_TokenEncryptionKeyContainer", "B2C_1A_SamlIdpCert"})

	resp, err := c.GraphServiceClient.TrustFramework().KeySets().Get(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, ks := range resp.GetValue() {
		id := ks.GetId()
		if a.Has(*id) {
			err = c.GraphServiceClient.TrustFramework().KeySets().ByTrustFrameworkKeySetId(*id).Delete(context.Background(), nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
