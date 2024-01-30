package cmd

import (
	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var deploy = &cobra.Command{
	Use:   "deploy [path to policies]",
	Short: "Deploy b2c policies.",
	Long:  `Deploy b2c policies to B2C identity experience framework.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())
		bd := internal.MustAbsPathFromFlag(cmd.Flags(), "build-dir")

		e.Deploy(bd)

		return nil
	},
}

func init() {
	globalFlags(deploy)
	deploy.Flags().StringP("build-dir", "b", "build", "Build directory")
	rootCmd.AddCommand(deploy)
}
