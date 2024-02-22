package cmd

import (
	"fmt"

	"com.schumann-it.go-ieftool/internal"
	"github.com/spf13/cobra"
)

var fix = &cobra.Command{
	Use:   "fix-app-registrations",
	Short: "Fix App Registrations",
	Long:  `Fix App Registration manifest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		err := e.FixAppRegistrations()
		if err != nil {
			return fmt.Errorf("errors occurred during fix app registrations process: \n%s", err.Error())
		}
		return err
	},
}

func init() {
	globalFlags(fix)
	rootCmd.AddCommand(fix)
}
