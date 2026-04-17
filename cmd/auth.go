package cmd

import (
	"context"
	"fmt"
	"os"

	"agrepl/pkg/auth"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication and identity",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to agrepl Cloud via GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		// In production, this would be a real Client ID for the Agrepl OAuth App
		clientID := os.Getenv("AGREPL_GITHUB_CLIENT_ID")
		if clientID == "" {
			// Placeholder for demonstration
			clientID = "Ov23lignS4D16X8p5U6K" 
		}

		fmt.Println("\033[36mInitiating GitHub Login...\033[0m")
		
		cfg, err := auth.PerformGitHubLogin(context.Background(), clientID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n\033[32m✓ Successfully logged in as %s (%s Tier)\033[0m\n", cfg.Email, cfg.Tier)
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from agrepl Cloud",
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.Logout(); err != nil {
			fmt.Printf("Already logged out or error: %v\n", err)
			return
		}
		fmt.Println("\033[32mSuccessfully logged out.\033[0m")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current auth status",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.LoadConfig()
		if err != nil {
			fmt.Println("\033[33mStatus: Not logged in (Local-only mode)\033[0m")
			return
		}
		fmt.Printf("\033[32mStatus: Logged in as %s\033[0m\n", cfg.Email)
		fmt.Printf("Tier: %s\n", cfg.Tier)
		fmt.Printf("Team ID: %s\n", cfg.TeamID)
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(authCmd)
}
