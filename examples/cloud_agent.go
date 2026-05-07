//go:build ignore

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"google.golang.org/api/idtoken"
)

// ANSI Color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

type CheckResult struct {
	Category string
	Provider string
	Status   string
	Details  string
}

func main() {
	fmt.Printf("%s%s=== Advanced Cloud Detection & Identity Agent ===%s\n", colorBold, colorCyan, colorReset)
	fmt.Println("This agent checks if you are RUNNING on a cloud or AUTHENTICATED to one.")
	fmt.Println("Note: 'aws s3 ls' works because you have local credentials, but you are likely")
	fmt.Println("not running on an EC2 instance (which explains why the metadata check fails).")
	fmt.Println(strings.Repeat("-", 85))

	results := make(chan CheckResult, 10)
	var wg sync.WaitGroup

	// --- Category 1: "Running On" (Metadata Services) ---
	wg.Add(3)
	go checkIMDS("AWS", "http://169.254.169.254/latest/meta-data/instance-id", nil, results, &wg)
	go checkIMDS("GCP", "http://metadata.google.internal/computeMetadata/v1/instance/id", map[string]string{"Metadata-Flavor": "Google"}, results, &wg)
	go checkIMDS("Azure", "http://169.254.169.254/metadata/instance/compute/vmId?api-version=2021-02-01", map[string]string{"Metadata": "true"}, results, &wg)

	// --- Category 2: "Authenticated To" (Public APIs) ---
	wg.Add(2)
	go checkAWSIdentity(results, &wg)
	go checkGCPIdentity(results, &wg)

	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Printf("%s%-15s | %-10s | %-12s | %s%s\n", colorBold, "Category", "Provider", "Status", "Details", colorReset)
	fmt.Println(strings.Repeat("-", 85))

	for res := range results {
		statusColor := colorRed
		if res.Status == "DETECTED" || res.Status == "AUTHED" {
			statusColor = colorGreen
		} else if res.Status == "TIMEOUT" || res.Status == "MISSING" {
			statusColor = colorYellow
		}

		fmt.Printf("%-15s | %-10s | %s%-12s%s | %s\n", 
			res.Category, res.Provider, statusColor, res.Status, colorReset, res.Details)
	}
	fmt.Println(strings.Repeat("-", 85))
	fmt.Printf("%sTIP:%s Use 'agrepl record -- go run examples/cloud_agent.go' to capture these API calls!%s\n", colorBold, colorCyan, colorReset)
}

func checkIMDS(name, url string, headers map[string]string, results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		results <- CheckResult{"Running On", name, "NOT FOUND", "Unreachable (Expected on local machine)"}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		results <- CheckResult{"Running On", name, "DETECTED", "Inside cloud environment"}
	} else {
		results <- CheckResult{"Running On", name, "FAILED", fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}
}

func checkAWSIdentity(results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Load AWS config (honors env vars and ~/.aws/credentials)
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		results <- CheckResult{"Auth To", "AWS", "ERROR", "Failed to load config"}
		return
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		results <- CheckResult{"Auth To", "AWS", "MISSING", "No valid AWS credentials found"}
		return
	}

	results <- CheckResult{"Auth To", "AWS", "AUTHED", fmt.Sprintf("Account: %s, Arn: %s", *identity.Account, *identity.Arn)}
}

func checkGCPIdentity(results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	// Simply check for Application Default Credentials
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		results <- CheckResult{"Auth To", "GCP", "MISSING", "GOOGLE_APPLICATION_CREDENTIALS not set"}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Try to validate credentials by creating an idtoken source
	_, err := idtoken.NewTokenSource(ctx, "https://example.com")
	if err != nil {
		results <- CheckResult{"Auth To", "GCP", "FAILED", "Invalid credentials"}
		return
	}

	results <- CheckResult{"Auth To", "GCP", "AUTHED", "GCP Credentials found and valid"}
}
