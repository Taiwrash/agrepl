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
			// Official agrepl OAuth App Client ID
			clientID = "Ov23lizFfb1GfxDvf8Ph" 
		}
		clientSecret := os.Getenv("AGREPL_GITHUB_CLIENT_SECRET")

		fmt.Println("\033[36mInitiating GitHub Login...\033[0m")
		
		cfg, err := auth.PerformGitHubLogin(context.Background(), clientID, clientSecret)
		if err != nil {
			if err.Error() == "failed to request device code: device flow not supported" {
				fmt.Fprintf(os.Stderr, "\n\033[31mError: Device Flow is not enabled for this GitHub OAuth App.\033[0m\n")
				fmt.Fprintf(os.Stderr, "To fix this:\n")
				fmt.Fprintf(os.Stderr, "1. Go to your GitHub App settings (Developer Settings > OAuth Apps).\n")
				fmt.Fprintf(os.Stderr, "2. Check the box \033[1m'Enable Device Flow'\033[0m.\n")
				fmt.Fprintf(os.Stderr, "3. Save changes and try again.\n\n")
				fmt.Fprintf(os.Stderr, "Current Client ID: %s\n", clientID)
			} else {
				fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			}
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

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade your account to the Pro tier",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := auth.LoadConfig()
		if err != nil {
			fmt.Println("\033[31mError: You must be logged in to upgrade. Run 'agrepl auth login'\033[0m")
			return
		}

		if cfg.Tier == "pro" || cfg.Tier == "enterprise" {
			fmt.Printf("Your account is already on the \033[1m%s\033[0m tier.\n", cfg.Tier)
			return
		}

		fmt.Println("\033[36mOpening checkout page at https://agrepl.dev/pricing/checkout...\033[0m")
		fmt.Println("Simulating payment success...")
		
		// In a real app, we would wait for a webhook or poll the backend.
		// For the demo, we just update the local config.
		cfg.Tier = "pro"
		if err := auth.SaveConfig(cfg); err != nil {
			fmt.Printf("\033[31mError saving updated tier: %v\033[0m\n", err)
			return
		}

		fmt.Printf("\n\033[32m✓ Successfully upgraded to Pro Tier!\033[0m\n")
		fmt.Println("You now have access to Team Sync, Semantic Diff, and more.")
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(authCmd)
}
