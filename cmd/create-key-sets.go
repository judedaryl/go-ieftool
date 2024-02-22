package cmd

import (
	"fmt"

	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var createKeySets = &cobra.Command{
	Use:   "create-key-sets",
	Short: "Create Key Sets",
	Long:  `Create Key Sets for Policies.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		err := e.CreateKeySets()
		if err != nil {
			return fmt.Errorf("errors occurred during key set creation process: \n%s", err.Error())
		}

		return nil
	},
}

func init() {
	globalFlags(createKeySets)
	rootCmd.AddCommand(createKeySets)
}
