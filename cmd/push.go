package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"agrepl/pkg/storage"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <run-id>",
	Short: "Uploads a run to remote storage",
	Long:  `The push command uploads a locally recorded agent execution to the configured R2 bucket.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runID := args[0]

		// Load local file directly as bytes to avoid unnecessary unmarshal/marshal
		filePath := filepath.Join(".agent-replay", "runs", fmt.Sprintf("%s.json", runID))
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading local run %s: %v\033[0m\n", runID, err)
			os.Exit(1)
		}

		rs, err := storage.NewRemoteStorage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError initializing remote storage: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[36mPushing run %s to remote storage...\033[0m\n", runID)
		if err := rs.Push(context.Background(), runID, data); err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[32mSuccessfully pushed run %s\033[0m\n", runID)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
