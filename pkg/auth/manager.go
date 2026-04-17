package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cli/oauth"
)

type Config struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	TeamID      string `json:"team_id"`
	Tier        string `json:"tier"` // "local", "pro", "enterprise"
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".agrepl")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "auth.json"), nil
}

func SaveConfig(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("not logged in. Run 'agrepl auth login'")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func Logout() error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func PerformGitHubLogin(ctx context.Context, clientID, clientSecret string) (*Config, error) {
	scopes := []string{"read:user", "user:email"}

	flow := &oauth.Flow{
		Host:         oauth.GitHubHost("https://github.com"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		HTTPClient:   http.DefaultClient,
		// If CallbackURI is provided, it tries Local Server flow first
		CallbackURI: "http://127.0.0.1/callback",
		// DisplayCode is called if it falls back to Device Flow
		DisplayCode: func(code, uri string) error {
			fmt.Printf("\n\033[33m! Action Required\033[0m\n")
			fmt.Printf("1. Copy your one-time code: \033[1;36m%s\033[0m\n", code)
			fmt.Printf("2. Open your browser to: \033[1;4m%s\033[0m\n", uri)
			fmt.Printf("\nWaiting for authorization...\n")
			return nil
		},
	}

	fmt.Println("\033[36mInitiating GitHub authorization...\033[0m")
	accessToken, err := flow.DetectFlow()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return finalizeLogin(accessToken.Token)
}

func finalizeLogin(token string) (*Config, error) {
	// In a real implementation, we would now call GitHub API to get user info
	cfg := &Config{
		AccessToken: token,
		Email:       "user@example.com", // Generic mock email
		TeamID:      "personal",
		Tier:        "local", // Default to Local (Free) tier
	}

	if err := SaveConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func IsFeatureAllowed(feature string) (bool, string) {
	cfg, err := LoadConfig()
	if err != nil {
		// Not logged in, only "local" features allowed
		if feature == "local_record" || feature == "local_replay" {
			return true, "local"
		}
		return false, "local"
	}

	switch cfg.Tier {
	case "enterprise":
		// Enterprise has everything
		return true, "enterprise"
	case "pro":
		// Pro has everything except enterprise-only features
		enterpriseOnly := map[string]bool{
			"guardrails": true,
			"audit_logs": true,
			"sso":        true,
		}
		if enterpriseOnly[feature] {
			return false, "pro"
		}
		return true, "pro"
	case "local":
		fallthrough
	default:
		// Local only has local features
		localFeatures := map[string]bool{
			"local_record": true,
			"local_replay": true,
		}
		if localFeatures[feature] {
			return true, "local"
		}
		return false, "local"
	}
}
