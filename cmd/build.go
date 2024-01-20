package cmd

import (
	"path"
	"path/filepath"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"

	logging "gopkg.in/op/go-logging.v1"
)

var build = &cobra.Command{
	Use:   "build",
	Short: "Build",
	Long:  `Build source policies and replacing template variables for given environments.`,
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logging.SetLevel(logging.INFO, "")
	},
	Run: func(cmd *cobra.Command, args []string) {
		cf, _ := cmd.Flags().GetString("config")
		en, _ := cmd.Flags().GetString("environment")
		sd, _ := cmd.Flags().GetString("source")
		if !filepath.IsAbs(sd) {
			sd, _ = filepath.Abs(sd)
		}
		dd, _ := cmd.Flags().GetString("destination")
		if !filepath.IsAbs(dd) {
			p := path.Join(sd, "..", dd)
			dd, _ = filepath.Abs(p)
		}

		e := internal.NewEnvironmentsFromConfig(cf, en)
		e.Build(sd, dd)
	},
}

func init() {
	build.Flags().StringP("config", "c", "config.yaml", "Path to the ieftool configuration file")
	build.Flags().StringP("source", "s", "source", "Source directory (current dir if omitted)")
	build.Flags().StringP("destination", "d", "build", "Destination directory (defaults to build relative to source)")
	build.Flags().StringP("environment", "e", "", "Environment to build (build all environments if omitted)")
	rootCmd.AddCommand(build)
}
