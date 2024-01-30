package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	logging "gopkg.in/op/go-logging.v1"
)

var rootCmd = &cobra.Command{
	Use:   "ieftool",
	Short: "Tooling for Azure B2C Identity Experience Framework",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		logging.SetLevel(logging.INFO, "")
	},
}

var completion = &cobra.Command{
	Use:                "completion [bash|zsh|fish|powershell]",
	Short:              "Generate completion script",
	Long:               `To load completions.`,
	DisableFlagParsing: true,
	ValidArgs:          []string{"bash", "zsh", "fish", "powershell"},
	Args:               cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completion)
}

func globalFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("config", "c", "", "Path to the ieftool configuration file (yaml)")
	cmd.Flags().StringP("environment", "e", "", "Environment to deploy (deploy all environments if omitted)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
