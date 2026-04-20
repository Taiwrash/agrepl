package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"agrepl/pkg/core"
	"agrepl/pkg/interceptor"
	"agrepl/pkg/proxy"
	"agrepl/pkg/storage"

	"github.com/spf13/cobra"
)

var recordCmd = &cobra.Command{
	Use:   "record <command> [args...]",
	Short: "Records an agent's execution",
	Long: `The record command executes a user-provided command, intercepting and recording
LLM calls and HTTP requests/responses. The execution trace is stored locally.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Record command called")

		// Initialize storage
		s, err := storage.NewJSONStorage(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError initializing storage: %v\033[0m\n", err)
			os.Exit(1)
		}

		runID, err := s.GetNextRunID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError getting next run ID: %v\033[0m\n", err)
			os.Exit(1)
		}

		currentRun := &core.Run{
			RunID: runID,
			Steps: []core.Step{}, // Initialize with an empty slice
		}

		// Start HTTP Proxy
		rand.Seed(time.Now().UnixNano())
		proxyPort := 8000 + rand.Intn(1000) // Random port between 8000 and 8999
		proxyAddr := "127.0.0.1:" + strconv.Itoa(proxyPort)

		httpInterceptor := interceptor.NewHTTPInterceptor(interceptor.ModeRecord, s, currentRun)
		httpProxy := proxy.NewHTTPProxy(proxyAddr, httpInterceptor)

		httpProxy.Start()
		defer func() {
			if err := httpProxy.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError stopping HTTP proxy: %v\033[0m\n", err)
			}
		}()

		// Prepare command to execute
		commandToExecute := args[0]
		commandArgs := args[1:]

		fmt.Printf("\033[36m[RECORD] Executing command: %s %v\033[0m\n", commandToExecute, commandArgs)

		execCmd := exec.Command(commandToExecute, commandArgs...)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		// Set proxy environment variables for the executed command
		execCmd.Env = os.Environ()
		proxyURL := "http://" + proxyAddr
		execCmd.Env = append(execCmd.Env, "HTTP_PROXY="+proxyURL)
		execCmd.Env = append(execCmd.Env, "HTTPS_PROXY="+proxyURL)
		execCmd.Env = append(execCmd.Env, "http_proxy="+proxyURL)
		execCmd.Env = append(execCmd.Env, "https_proxy="+proxyURL)
		execCmd.Env = append(execCmd.Env, "ALL_PROXY="+proxyURL)
		execCmd.Env = append(execCmd.Env, "all_proxy="+proxyURL)
		execCmd.Env = append(execCmd.Env, "NO_PROXY=localhost,127.0.0.1")
		execCmd.Env = append(execCmd.Env, "no_proxy=localhost,127.0.0.1")

		// Trust the local Root CA
		caPath, _ := filepath.Abs(filepath.Join(".agent-replay", "ca", "ca.crt"))
		execCmd.Env = append(execCmd.Env, "SSL_CERT_FILE="+caPath)
		execCmd.Env = append(execCmd.Env, "REQUESTS_CA_BUNDLE="+caPath) // For Python requests library
		execCmd.Env = append(execCmd.Env, "NODE_EXTRA_CA_CERTS="+caPath) // For Node.js

		// Set agrepl mode for LLM interception
		execCmd.Env = append(execCmd.Env, "AGREPL_MODE=record")
		execCmd.Env = append(execCmd.Env, "AGREPL_RUN_ID="+runID)

		// Attempt to save the run, even if the command failed or was interrupted.
		// This ensures partial runs are saved.
		var saveErr error
		if saveErr = s.SaveRun(currentRun); saveErr != nil { // Assign to existing saveErr
			fmt.Fprintf(os.Stderr, "\033[31mError saving recorded run (possibly partial): %v\033[0m\n", saveErr)
		}

		// Retrieve the run ID. It might be set by LLM interception or be the initial runID.
		// This must be done before db.SaveMetadata and the final print statement.
		finalRunID := os.Getenv("AGREPL_CURRENT_RUN_ID")
		if finalRunID == "" {
			finalRunID = runID // Fallback to the initially generated ID
		}
		currentRun.RunID = finalRunID // Update currentRun's ID to the final one

		// Attempt to update SQLite index with metadata for the run (partial or complete)
		if db, dbErr := storage.NewDB("."); dbErr == nil {
			defer db.Close()
			status := "completed"
			commandExecuted := commandToExecute + " " + strings.Join(commandArgs, " ")
			if err != nil { // If execCmd.Run() had an error
				status = "failed"
			}
			db.SaveMetadata(&storage.RunMetadata{
				RunID:      finalRunID, // Use finalRunID here
				Command:    commandExecuted,
				CreatedAt:  time.Now(),
				TotalSteps: len(currentRun.Steps),
				Status:     status,
			})
		} else {
			fmt.Fprintf(os.Stderr, "\033[31mError initializing DB for metadata update: %v\033[0m\n", dbErr)
		}

		// Now handle the error from execCmd.Run() if it occurred
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError executing command: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[32m[RECORD] Recorded run with ID: %s. Total steps: %d\033[0m\n", finalRunID, len(currentRun.Steps))
	},
}

func init() {
	rootCmd.AddCommand(recordCmd)
}
