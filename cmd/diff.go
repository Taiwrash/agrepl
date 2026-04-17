package cmd

import (
	"bytes"
	"fmt"
	"os"

	"agrepl/pkg/core"
	"agrepl/pkg/storage"

	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <run-id-1> <run-id-2>",
	Short: "Compares two agent runs and shows differences",
	Long:  `The diff command compares two recorded agent executions, highlighting differences in steps, status codes, and responses.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runID1 := args[0]
		runID2 := args[1]

		s, err := storage.NewJSONStorage(".")
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError: %v\033[0m\n", err)
			os.Exit(1)
		}

		run1, err := s.LoadRun(runID1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading %s: %v\033[0m\n", runID1, err)
			os.Exit(1)
		}

		run2, err := s.LoadRun(runID2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mError loading %s: %v\033[0m\n", runID2, err)
			os.Exit(1)
		}

		fmt.Printf("\033[36mComparing %s and %s\033[0m\n\n", runID1, runID2)

		maxSteps := len(run1.Steps)
		if len(run2.Steps) > maxSteps {
			maxSteps = len(run2.Steps)
		}

		diffFound := false
		for i := 0; i < maxSteps; i++ {
			if i >= len(run1.Steps) {
				fmt.Printf("\033[32m[+ STEP %d] (Only in %s): %s\033[0m\n", i, runID2, stepSummary(run2.Steps[i]))
				diffFound = true
				continue
			}
			if i >= len(run2.Steps) {
				fmt.Printf("\033[31m[- STEP %d] (Only in %s): %s\033[0m\n", i, runID1, stepSummary(run1.Steps[i]))
				diffFound = true
				continue
			}

			s1 := run1.Steps[i]
			s2 := run2.Steps[i]

			if s1.Type != s2.Type {
				fmt.Printf("\033[33m[Δ STEP %d] Type Mismatch: %s -> %s\033[0m\n", i, s1.Type, s2.Type)
				diffFound = true
				continue
			}

			switch s1.Type {
			case core.StepTypeHTTP:
				diffHTTP(i, s1, s2)
			case core.StepTypeLLM:
				diffLLM(i, s1, s2)
			}
		}

		if !diffFound {
			fmt.Println("\033[32mNo differences found between runs.\033[0m")
		}
	},
}

func stepSummary(s core.Step) string {
	switch s.Type {
	case core.StepTypeHTTP:
		return fmt.Sprintf("HTTP %s %s", s.Request.Method, s.Request.URL)
	case core.StepTypeLLM:
		return fmt.Sprintf("LLM %s", s.LLMInput.Model)
	default:
		return string(s.Type)
	}
}

func diffHTTP(idx int, s1, s2 core.Step) {
	headerPrinted := false
	printHeader := func() {
		if !headerPrinted {
			fmt.Printf("\033[33m[Δ STEP %d] HTTP %s %s\033[0m\n", idx, s1.Request.Method, s1.Request.URL)
			headerPrinted = true
		}
	}

	if s1.Request.URL != s2.Request.URL {
		printHeader()
		fmt.Printf("  URL: \033[31m%s\033[0m -> \033[32m%s\033[0m\n", s1.Request.URL, s2.Request.URL)
	}

	if s1.Response.StatusCode != s2.Response.StatusCode {
		printHeader()
		fmt.Printf("  Status: \033[31m%d\033[0m -> \033[32m%d\033[0m\n", s1.Response.StatusCode, s2.Response.StatusCode)
	}

	if !bytes.Equal(s1.Response.Body, s2.Response.Body) {
		printHeader()
		fmt.Printf("  Body Differs:\n")
		fmt.Printf("    \033[31m- %s\033[0m\n", truncate(string(s1.Response.Body)))
		fmt.Printf("    \033[32m+ %s\033[0m\n", truncate(string(s2.Response.Body)))
	}
}

func diffLLM(idx int, s1, s2 core.Step) {
	headerPrinted := false
	printHeader := func() {
		if !headerPrinted {
			fmt.Printf("\033[33m[Δ STEP %d] LLM %s\033[0m\n", idx, s1.LLMInput.Model)
			headerPrinted = true
		}
	}

	if s1.LLMInput.Model != s2.LLMInput.Model {
		printHeader()
		fmt.Printf("  Model: \033[31m%s\033[0m -> \033[32m%s\033[0m\n", s1.LLMInput.Model, s2.LLMInput.Model)
	}

	// Compare responses
	r1 := ""
	r2 := ""
	if len(s1.LLMOutput.Response.Candidates) > 0 && len(s1.LLMOutput.Response.Candidates[0].Content.Parts) > 0 {
		r1 = s1.LLMOutput.Response.Candidates[0].Content.Parts[0].Text
	}
	if len(s2.LLMOutput.Response.Candidates) > 0 && len(s2.LLMOutput.Response.Candidates[0].Content.Parts) > 0 {
		r2 = s2.LLMOutput.Response.Candidates[0].Content.Parts[0].Text
	}

	if r1 != r2 {
		printHeader()
		fmt.Printf("  Response: \033[31m%s\033[0m -> \033[32m%s\033[0m\n", truncate(r1), truncate(r2))
	}
}

func truncate(s string) string {
	if len(s) > 50 {
		return s[:47] + "..."
	}
	return s
}

func init() {
	rootCmd.AddCommand(diffCmd)
}
