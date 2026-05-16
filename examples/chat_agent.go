package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"agrepl/pkg/llmclient"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Disable TLS verification if we are being intercepted by agrepl
	if os.Getenv("AGREPL_MODE") != "" {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// Create the intercepted model
	modelName := os.Getenv("GEMINI_MODEL")
	if modelName == "" {
		modelName = "gemini-2.0-flash"
	}
	model, err := llmclient.NewGenerativeModel(ctx, apiKey, modelName)
	if err != nil {
		log.Fatalf("Error creating model: %v", err)
	}

	fmt.Println("--- Gemini Chat Agent (Intercepted) ---")
	fmt.Println("Type 'exit' or 'quit' to stop.")
	fmt.Println("---------------------------------------")

	scanner := bufio.NewScanner(os.Stdin)
	var history []*genai.Content

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if strings.ToLower(input) == "exit" || strings.ToLower(input) == "quit" {
			break
		}

		// Add user input to history
		history = append(history, &genai.Content{
			Role: "user",
			Parts: []*genai.Part{
				{Text: input},
			},
		})

		// Generate response
		resp, err := model.GenerateContent(ctx, history, nil)
		if err != nil {
			log.Printf("Error generating content: %v", err)
			continue
		}

		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			responseText := resp.Candidates[0].Content.Parts[0].Text
			fmt.Printf("Gemini: %s\n", responseText)

			// Add model response to history
			history = append(history, &genai.Content{
				Role: "model",
				Parts: []*genai.Part{
					{Text: responseText},
				},
			})
		} else {
			fmt.Println("Gemini: [No response]")
		}
	}

	fmt.Println("Goodbye!")
}
