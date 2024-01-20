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
	Run: func(cmd *cobra.Command, args []string) {
		cf, _ := cmd.Flags().GetString("config")
		en, _ := cmd.Flags().GetString("environment")
		e := internal.NewEnvironmentsFromConfig(cf, en)
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
	},
}

func init() {
	list.Flags().StringP("config", "c", "ieftool.config", "Path to the ieftool configuration file (yaml)")
	list.Flags().StringP("environment", "e", "", "Environment to deploy (deploy all environments if omitted)")
	rootCmd.AddCommand(list)
}
