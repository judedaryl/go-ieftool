package cmd

import (
	"log"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var fix = &cobra.Command{
	Use:   "fix-app-registrations",
	Short: "Fix App Registrations",
	Long:  `Fix App Registration manifest.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tid, err := cmd.Flags().GetString("tenant-id")
		if err != nil {
			log.Fatalf("could not parse flag 'tenant-id': \n%s", err.Error())
		}
		cid, err := cmd.Flags().GetString("client-id")
		if err != nil {
			log.Fatalf("could not parse flag 'client-id': \n%s", err.Error())
		}
		aid, err := cmd.Flags().GetString("application-object-id")
		if err != nil {
			log.Fatalf("could not parse flag 'application-id': \n%s", err.Error())
		}

		g, err := internal.NewGraphClient(tid, cid)
		if err != nil {
			log.Fatalf("could not create graph client: \n%s", err.Error())
		}

		err = g.FixAppRegistration(aid)
		return err
	},
}

func init() {
	fix.Flags().String("tenant-id", "", "Tenant ID")
	fix.Flags().String("client-id", "", "Client ID")
	fix.Flags().String("application-object-id", "", "Application Object ID")

	rootCmd.AddCommand(fix)
}
