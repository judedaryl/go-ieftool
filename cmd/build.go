package cmd

import (
	"os"

	"com.go.ieftool/internal"
	"github.com/spf13/cobra"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	logging "gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

var build = &cobra.Command{
	Use:   "build [path to source code]",
	Short: "Build source policies and replacing variables.",
	Long:  `Build source policies and replacing template variables with their corresponding values.`,
	Args:  cobra.MinimumNArgs(1),
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		yqlib.InitExpressionParser()
		logging.SetLevel(logging.INFO, "")
	},
	Run: func(cmd *cobra.Command, args []string) {
		configurationPath, err := cmd.Flags().GetString("config")
		internal.Check(err)
		config, err := os.ReadFile(configurationPath)
		internal.Check(err)

		yml := make(map[string]interface{})
		err = yaml.Unmarshal(config, yml)
		internal.Check(err)
		internal.Build(configurationPath, args[0], "", args[1])
	},
}

func init() {
	build.Flags().StringP("config", "c", "ieftool.config", "Path to the ieftool configuration file (yaml)")
	rootCmd.AddCommand(build)
}
