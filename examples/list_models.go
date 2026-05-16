//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	fmt.Println("--- Available Gemini Models ---")
	
	// ListModels uses a pager
	it := client.Models.List(ctx, nil)
	for {
		model, err := it.Next()
		if err != nil {
			// it.Next() returns io.EOF when done
			break
		}
		fmt.Printf("Model: %-30s | Methods: %v\n", model.Name, model.SupportedGenerationMethods)
	}
}
