package cmd

import (
	"log"
	"os"
	"sync"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var deploy = &cobra.Command{
	Use:   "deploy [path to policies]",
	Short: "Deploy b2c policies.",
	Long:  `Deploy b2c policies to B2C identity experience framework.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		tenantId := os.Getenv("B2C_TENANT_ID")
		clientId := os.Getenv("B2C_CLIENT_ID")
		clientSecret := os.Getenv("B2C_CLIENT_SECRET")

		if tenantId == "" {
			log.Fatalln("Environment variable B2C_TENANT_ID has not been set.")
		}

		if clientId == "" {
			log.Fatalln("Environment variable B2C_CLIENT_ID has not been set.")
		}

		if clientSecret == "" {
			log.Fatalln("Environment variable B2C_CLIENT_SECRET has not been set.")
		}

		filePath := args[0]
		token := internal.GetToken(clientId, clientSecret, tenantId)

		policies := []internal.Policy{}
		policies = internal.GetPolicies(filePath, policies)
		batchedPolicies := internal.CreateBatchedArray(policies)

		for i, batch := range batchedPolicies {
			log.Printf("Processing batch %d", i)
			uploadPolicies(token.AccessToken, batch)
		}
	},
}

func init() {
	rootCmd.AddCommand(deploy)
}

func uploadPolicies(token string, policies []internal.Policy) {
	var wg sync.WaitGroup
	wg.Add(len(policies))

	for _, p := range policies {
		go uploadPolicy(token, p, &wg)
	}
	wg.Wait()
}

func uploadPolicy(token string, policy internal.Policy, wg *sync.WaitGroup) {
	defer wg.Done()

	content, err := os.ReadFile(policy.Path)
	internal.Check(err)

	internal.UploadPolicy(token, policy.PolicyId, content)
	log.Printf("\tUploaded: %s\n", policy.PolicyId)
}
