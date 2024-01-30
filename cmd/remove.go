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
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		err := e.DeleteRemotePolicies()
		if err != nil {
			log.Fatalf("Failed to remove policies %s", err.Error())
		}

		return nil
	},
}

func init() {
	globalFlags(remove)
	rootCmd.AddCommand(remove)
}
