package cmd

import (
	"path/filepath"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var deploy = &cobra.Command{
	Use:   "deploy [path to policies]",
	Short: "Deploy b2c policies.",
	Long:  `Deploy b2c policies to B2C identity experience framework.`,
	Run: func(cmd *cobra.Command, args []string) {
		cf, _ := cmd.Flags().GetString("config")
		en, _ := cmd.Flags().GetString("environment")
		bd, _ := cmd.Flags().GetString("build-dir")
		if !filepath.IsAbs(bd) {
			bd, _ = filepath.Abs(bd)
		}
		e := internal.NewEnvironmentsFromConfig(cf, en)
		e.Deploy(bd)
	},
}

func init() {
	deploy.Flags().StringP("config", "c", "ieftool.config", "Path to the ieftool configuration file (yaml)")
	deploy.Flags().StringP("environment", "e", "", "Environment to deploy (deploy all environments if omitted)")
	deploy.Flags().StringP("build-dir", "b", "", "Build directory")
	rootCmd.AddCommand(deploy)
}
