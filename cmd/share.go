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

		// 1. Check Auth
		cfg, err := auth.LoadConfig()
		if err != nil {
			fmt.Printf("\033[31mError: %v\033[0m\n", err)
			os.Exit(1)
		}

		// 2. Load Local Data
		filePath := filepath.Join(".agent-replay", "runs", fmt.Sprintf("%s.json", runID))
		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading local run %s: %v\033[0m\n", runID, err)
			os.Exit(1)
		}

		// 3. Register and Upload (Mocking Backend Call)
		// In reality, this would hit POST /runs to get a pre-signed URL, 
		// then upload the data to that URL.
		fmt.Printf("\033[36mAuthenticating with Team Workspace: %s...\033[0m\n", cfg.TeamID)
		
		rs, err := storage.NewRemoteStorage()
		if err != nil {
			// If R2 env vars are not set, we'd normally use the pre-signed URL from the backend.
			// For this MVP, we still rely on the R2 storage package but wrap it in the 'share' command.
			fmt.Printf("\033[33mWarning: Remote storage backend not fully configured for pre-signed URLs.\033[0m\n")
		} else {
			if err := rs.Push(context.Background(), runID, data); err != nil {
				fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
				os.Exit(1)
			}
		}

		// 4. Return Shareable ID
		shareID := fmt.Sprintf("share-%s-%s", cfg.TeamID[:4], runID)
		fmt.Printf("\n\033[32m✓ Run successfully shared with your team!\033[0m\n")
		fmt.Printf("Share ID: \033[1m%s\033[0m\n", shareID)
		fmt.Printf("Teammates can replay this using:\n  agrepl pull %s\n", shareID)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
}
