package cmd

import (
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
	Short: "Log in to agrepl Cloud",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Attempting to log in to agrepl Cloud...")
		
		// For now, we simulate a login success for development
		// In a real implementation, this would initiate the Device Flow
		mockCfg := &auth.Config{
			AccessToken: "mock-jwt-token-pro-tier",
			Email:       "dev@example.com",
			TeamID:      "team-id-123",
		}

		if err := auth.SaveConfig(mockCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving auth config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[32mSuccessfully logged in as %s (Pro Tier)\033[0m\n", mockCfg.Email)
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
		fmt.Printf("Team ID: %s\n", cfg.TeamID)
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(authCmd)
}
