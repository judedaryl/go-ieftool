package cmd

import (
	"fmt"

	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var deleteKeySets = &cobra.Command{
	Use:   "delete-key-sets",
	Short: "Delete Key Sets",
	Long:  `Delete Key Sets for Policies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		err := e.DeleteKeySets()
		if err != nil {
			return fmt.Errorf("errors occurred during key set deletion process: \n%s", err.Error())
		}

		return nil
	},
}

func init() {
	globalFlags(deleteKeySets)
	rootCmd.AddCommand(deleteKeySets)
}
