package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"agrepl/pkg/auth"
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

		// Check if Team Sync is allowed for the current tier
		allowed, tier := auth.IsFeatureAllowed("team_sync")
		if !allowed {
			fmt.Printf("\033[33mThe 'push' command (Team Sync) is a Pro feature.\033[0m\n")
			fmt.Printf("Your current tier is: \033[1m%s\033[0m\n", tier)
			fmt.Println("Please upgrade to Pro at https://agrepl.dev/pricing to share runs with your team.")
			os.Exit(1)
		}

		// Load local file directly as bytes to avoid unnecessary unmarshal/marshal
		filePath := filepath.Join(".agent-replay", "runs", fmt.Sprintf("%s.json", runID))
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading local run %s: %v\033[0m\n", runID, err)
			os.Exit(1)
		}

		// Load auth config to get TeamID
		cfg, err := auth.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: You must be logged in to push. Run 'agrepl auth login'\033[0m\n")
			os.Exit(1)
		}

		rs, err := storage.NewRemoteStorage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError initializing remote storage: %v\033[0m\n", err)
			os.Exit(1)
		}
		rs.Prefix = cfg.TeamID

		teamPrefix := cfg.TeamID
		if len(teamPrefix) > 4 {
			teamPrefix = teamPrefix[:4]
		}
		remoteID := fmt.Sprintf("share-%s-%s", teamPrefix, runID)

		fmt.Printf("\033[36mPushing run %s to remote storage as %s (Team: %s)...\033[0m\n", runID, remoteID, cfg.TeamID)
		if err := rs.Push(context.Background(), remoteID, data); err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[32mSuccessfully pushed run %s. Available as: %s\033[0m\n", runID, remoteID)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
