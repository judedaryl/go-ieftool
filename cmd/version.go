package cmd

import (
	"com.schumann-it.go-ieftool/version"
	"github.com/spf13/cobra"
)

var printVersion = &cobra.Command{
	Use:   "version",
	Short: "Print version.",
	Long:  `Print version.`,
	Run: func(cmd *cobra.Command, args []string) {
		println(version.String())
	},
}

func init() {
	rootCmd.AddCommand(printVersion)
}
