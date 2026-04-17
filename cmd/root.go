package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "agrepl",
	Short: "agrepl is a lightweight CLI tool for recording and replaying AI agent executions",
	Long: `agrepl (agent-replay) is a lightweight CLI tool that enables developers to record and deterministically replay AI agent executions.

It focuses on reproducibility, debugging, and offline testing of agent workflows by capturing LLM interactions and external tool/API calls.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action if no subcommand is given
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.agrepl.yaml)")
	// Cobra also supports local flags, which will only run when this command
	// is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
