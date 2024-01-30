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

		errs := e.DeleteRemotePolicies()
		if errs != nil {
			for _, err := range errs {
				log.Println(err)
			}
		}

		return nil
	},
}

func init() {
	globalFlags(remove)
	rootCmd.AddCommand(remove)
}
