package cmd

import (
	"log"

	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

var rootCmd = &cobra.Command{
	Use:   "ieftool",
	Short: "Tooling for Azure B2C Identity Experience Framework",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logging.SetLevel(logging.INFO, "")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
