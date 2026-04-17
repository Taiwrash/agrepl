package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"agrepl/pkg/core"
	"agrepl/pkg/storage"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <run-id>",
	Short: "Downloads a run from remote storage",
	Long:  `The pull command downloads an agent execution from the configured R2 bucket and stores it locally.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runID := args[0]

		rs, err := storage.NewRemoteStorage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError initializing remote storage: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[36mPulling run %s from remote storage...\033[0m\n", runID)
		data, err := rs.Pull(context.Background(), runID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			os.Exit(1)
		}

		// Save locally
		localDir := filepath.Join(".agent-replay", "runs")
		os.MkdirAll(localDir, 0755)
		filePath := filepath.Join(localDir, fmt.Sprintf("%s.json", runID))
		if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError saving run locally: %v\033[0m\n", err)
			os.Exit(1)
		}

		// Update SQLite index
		var run core.Run
		if err := json.Unmarshal(data, &run); err == nil {
			db, err := storage.NewDB(".")
			if err == nil {
				defer db.Close()
				db.SaveMetadata(&storage.RunMetadata{
					RunID:      run.RunID,
					Command:    "pulled-from-remote", // We don't know the original command unless we store it in the JSON
					CreatedAt:  time.Now(),
					TotalSteps: len(run.Steps),
					Status:     "completed",
				})
			}
		}

		fmt.Printf("\033[32mSuccessfully pulled run %s\033[0m\n", runID)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
