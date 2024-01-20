package cmd

import (
	"log"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var remove = &cobra.Command{
	Use:   "remove",
	Short: "Delete remote b2c policies.",
	Long:  `Delete remote b2c policies from B2C identity experience framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		cf, _ := cmd.Flags().GetString("config")
		en, _ := cmd.Flags().GetString("environment")
		e := internal.NewEnvironmentsFromConfig(cf, en)
		errs := e.DeleteRemotePolicies()
		if errs != nil {
			for _, err := range errs {
				log.Println(err)
			}
		}
	},
}

func init() {
	remove.Flags().StringP("config", "c", "ieftool.config", "Path to the ieftool configuration file (yaml)")
	remove.Flags().StringP("environment", "e", "", "Environment to deploy (deploy all environments if omitted)")
	rootCmd.AddCommand(remove)
}
