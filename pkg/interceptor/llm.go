package interceptor

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"agrepl/pkg/core"
	"agrepl/pkg/storage"

	"google.golang.org/genai"
)

// LLMInterceptor wraps a genai.Client to intercept LLM calls for a specific model.
type LLMInterceptor struct {
	Client      *genai.Client
	ModelName   string
	Mode        Mode
	Storage     storage.Storage
	CurrentRun  *core.Run
	usedSteps   map[int]bool // Track consumed steps
	mu          sync.Mutex
}

// NewLLMInterceptor creates a new LLMInterceptor.
func NewLLMInterceptor(client *genai.Client, modelName string, mode Mode, s storage.Storage, currentRun *core.Run) *LLMInterceptor {
	return &LLMInterceptor{
		Client:      client,
		ModelName:   modelName,
		Mode:        mode,
		Storage:     s,
		CurrentRun:  currentRun,
		usedSteps:   make(map[int]bool),
	}
}

// GenerateContent intercepts the genai.Models.GenerateContent call.
func (i *LLMInterceptor) GenerateContent(ctx context.Context, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	switch i.Mode {
	case ModeRecord:
		return i.recordGenerateContent(ctx, contents, cfg)
	case ModeReplay:
		return i.replayGenerateContent(ctx, contents, cfg)
	case ModePassthrough:
		fallthrough
	default:
		return i.Client.Models.GenerateContent(ctx, i.ModelName, contents, cfg)
	}
}

func (i *LLMInterceptor) recordGenerateContent(ctx context.Context, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	resp, err := i.Client.Models.GenerateContent(ctx, i.ModelName, contents, cfg)
	if err != nil {
		return nil, err
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	llmStep := core.Step{
		Type: core.StepTypeLLM,
		LLMInput: &core.GeminiGenerateContentRequest{
			Model:    i.ModelName,
			Contents: contents,
			Config:   cfg,
		},
		LLMOutput: &core.GeminiGenerateContentResponse{
			Response: resp,
		},
	}
	i.CurrentRun.Steps = append(i.CurrentRun.Steps, llmStep)
	if i.Storage != nil {
		i.Storage.AppendStep(i.CurrentRun.RunID, llmStep)
	}
	fmt.Printf("%s[RECORD]%s Captured LLM: %s\n", colorCyan, colorReset, i.ModelName)

	return resp, nil
}

func (i *LLMInterceptor) replayGenerateContent(ctx context.Context, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Deep copy incoming contents for matching
	incomingContentsCopy := make([]*genai.Content, len(contents))
	for idx, c := range contents {
		contentCopy := *c
		incomingContentsCopy[idx] = &contentCopy
	}

	for idx, step := range i.CurrentRun.Steps {
		if i.usedSteps[idx] || step.Type != core.StepTypeLLM {
			continue
		}

		recordedReq := step.LLMInput
		if recordedReq == nil {
			continue
		}

		// Exact matching for LLM requests
		modelMatch := recordedReq.Model == i.ModelName
		contentsMatch := compareContents(recordedReq.Contents, incomingContentsCopy)

		if modelMatch && contentsMatch {
			i.usedSteps[idx] = true
			fmt.Printf("%s[REPLAY] Matched LLM call: %s%s\n", colorGreen, i.ModelName, colorReset)
			fmt.Printf("         Returning recorded Gemini response (Run: %s, Step: %d)\n", i.CurrentRun.RunID, idx)
			return step.LLMOutput.Response, nil
		}
	}

	fmt.Printf("%s[REPLAY] ERROR: No matching recorded LLM step found for model %s%s\n", colorRed, i.ModelName, colorReset)
	return nil, fmt.Errorf("replay error: no matching recorded LLM step found for model %s", i.ModelName)
}

func compareContents(a, b []*genai.Content) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
