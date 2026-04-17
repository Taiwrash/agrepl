package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"agrepl/pkg/interceptor"
	"agrepl/pkg/proxy"
	"agrepl/pkg/storage"

	"github.com/spf13/cobra"
)

var fallback bool

var replayCmd = &cobra.Command{
	Use:   "replay <run-id> <command> [args...]",
	Short: "Replays a previously recorded agent execution",
	Long: `The replay command re-executes a previously recorded agent run,
returning recorded responses instead of making real calls, ensuring determinism.`,
	Args: cobra.MinimumNArgs(2), // run-id + command
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Replay command called")
		runID := args[0]

		// Initialize storage
		s, err := storage.NewJSONStorage(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError initializing storage: %v\033[0m\n", err)
			os.Exit(1)
		}

		run, err := s.LoadRun(runID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading run '%s': %v\033[0m\n", runID, err)
			os.Exit(1)
		}

		// Start HTTP Proxy in replay mode
		rand.Seed(time.Now().UnixNano())
		proxyPort := 9000 + rand.Intn(1000) // Random port between 9000 and 9999
		proxyAddr := "127.0.0.1:" + strconv.Itoa(proxyPort)

		httpInterceptor := interceptor.NewHTTPInterceptor(interceptor.ModeReplay, s, run)
		httpInterceptor.Fallback = fallback // Set fallback mode

		httpProxy := proxy.NewHTTPProxy(proxyAddr, httpInterceptor)

		httpProxy.Start()
		defer func() {
			if err := httpProxy.Stop(); err != nil {
				fmt.Fprintf(os.Stderr, "\033[31mError stopping HTTP proxy: %v\033[0m\n", err)
			}
		}()

		// Prepare command to execute
		commandToExecute := args[1]
		commandArgs := args[2:]

		fmt.Printf("\033[32m[REPLAY] Replaying run with ID: %s by executing command: %s %v\033[0m\n", runID, commandToExecute, commandArgs)

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

		// Set agrepl mode and run ID for LLM interception
		execCmd.Env = append(execCmd.Env, "AGREPL_MODE=replay")
		execCmd.Env = append(execCmd.Env, "AGREPL_RUN_ID="+runID)

		if err := execCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError executing command during replay: %v\033[0m\n", err)
			os.Exit(1)
		}

		fmt.Printf("\033[32m[REPLAY] Replay of run %s finished.\033[0m\n", runID)
	},
}

func init() {
	replayCmd.Flags().BoolVarP(&fallback, "fallback", "f", false, "Hit real network if recorded request is not found")
	rootCmd.AddCommand(replayCmd)
}
