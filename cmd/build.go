package cmd

import (
	"fmt"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"
)

var build = &cobra.Command{
	Use:   "build",
	Short: "Build",
	Long:  `Build source policies and replacing template variables for given environments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		e := internal.MustNewEnvironmentsFromFlags(cmd.Flags())
		sd := internal.MustAbsPathFromFlag(cmd.Flags(), "source")
		dd := internal.MustAbsPathFromFlag(cmd.Flags(), "destination")

		err := e.Build(sd, dd)
		if err != nil {
			return fmt.Errorf("errors occurred during build process: \n%s", err.Error())
		}

		return nil
	},
}

func init() {
	globalFlags(build)
	build.Flags().StringP("source", "s", "src", "Source directory")
	build.Flags().StringP("destination", "d", "build", "Destination directory")
	rootCmd.AddCommand(build)
}
