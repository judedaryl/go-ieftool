package cmd

import (
	"log"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:   "list [path to policies]",
	Short: "List remote b2c policies.",
	Long:  `List remote b2c policies from B2C identity experience framework.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())

		ps, err := e.ListRemotePolicies()
		if err != nil {
			log.Fatalf("Failed to llist policies %v", err)
		}

		for n, l := range ps {
			log.Println(n)
			for _, i := range l {
				log.Println(i)
			}
		}

		return nil
	},
}

func init() {
	globalFlags(list)
	rootCmd.AddCommand(list)
}
