package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"agrepl/pkg/llmclient"
	"google.golang.org/genai"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Create the intercepted model
	modelName := "gemini-2.0-flash"
	model, err := llmclient.NewGenerativeModel(ctx, apiKey, modelName)
	if err != nil {
		log.Fatalf("Error creating model: %v", err)
	}

	fmt.Printf("--- Calling Gemini (%s) ---\n", modelName)
	
	// genai.Text returns []*genai.Content
	contents := genai.Text("Explain the concept of 'determinism' in software testing in one sentence.")

	resp, err := model.GenerateContent(ctx, contents, nil)
	if err != nil {
		log.Fatalf("Error generating content: %v", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		fmt.Printf("Response: %s\n", resp.Candidates[0].Content.Parts[0].Text) // Added .Text for clarity if it's a Part
	} else {
		fmt.Println("No response generated.")
	}
}
