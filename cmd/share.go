package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"agrepl/pkg/auth"
	"agrepl/pkg/storage"

	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share <run-id>",
	Short: "Shares a run with your team and generates a link",
	Long:  `The share command uploads a run to the agrepl Cloud and provides a shareable ID or URL.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runID := args[0]

		// 1. Check Auth & Tier
		allowed, tier := auth.IsFeatureAllowed("team_sync")
		if !allowed {
			fmt.Printf("\033[33mThe 'share' command is a Pro feature.\033[0m\n")
			fmt.Printf("Your current tier is: \033[1m%s\033[0m\n", tier)
			fmt.Println("Please upgrade to Pro at https://agrepl.dev/pricing to share runs.")
			os.Exit(1)
		}

		cfg, err := auth.LoadConfig()
		if err != nil {
			fmt.Printf("\033[31mError loading config: %v\033[0m\n", err)
			os.Exit(1)
		}

		// 2. Load Local Data
		filePath := filepath.Join(".agent-replay", "runs", fmt.Sprintf("%s.json", runID))
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading local run %s: %v\033[0m\n", runID, err)
			os.Exit(1)
		}

		// 3. Register and Upload
		fmt.Printf("\033[36mAuthenticating with Team Workspace: %s...\033[0m\n", cfg.TeamID)

		rs, err := storage.NewRemoteStorage()
		if err != nil {
			fmt.Printf("\033[33mWarning: Remote storage not configured (Env vars missing).\033[0m\n")
			fmt.Println("To enable real uploads, set AGREPL_R2_* environment variables.")
		} else {
			fmt.Printf("\033[36mUploading run to cloud...\033[0m\n")
			if err := rs.Push(context.Background(), runID, data); err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mUpload failed: %v\033[0m\n", err)
				os.Exit(1)
			}
		}

		// 4. Return Shareable ID
		teamPrefix := cfg.TeamID
		if len(teamPrefix) > 4 {
			teamPrefix = teamPrefix[:4]
		}
		shareID := fmt.Sprintf("share-%s-%s", teamPrefix, runID)
		fmt.Printf("\n\033[32m✓ Run successfully shared!\033[0m\n")
		fmt.Printf("Share ID: \033[1m%s\033[0m\n", shareID)
		fmt.Printf("Teammates can replay this using:\n  agrepl pull %s\n", shareID)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
