package core

import (
	"net/http"

	"google.golang.org/genai"
)

// Run represents a single recorded execution of an agent.
type Run struct {
	RunID           string `json:"run_id"`
	OriginalCommand string `json:"original_command"`
	Steps           []Step `json:"steps"`
}

// StepType defines the type of an interaction step.
type StepType string

const (
	StepTypeLLM  StepType = "llm"
	StepTypeHTTP StepType = "http"
)

// Step represents a single interaction within an agent's execution.
// It can be either an LLM interaction or an HTTP interaction.
type Step struct {
	Type      StepType                      `json:"type"`
	LLMInput  *GeminiGenerateContentRequest `json:"llm_input,omitempty"`  // Gemini-specific LLM input structure
	LLMOutput *GeminiGenerateContentResponse `json:"llm_output,omitempty"` // Gemini-specific LLM output structure
	Request   *HTTPRequest                  `json:"request,omitempty"`    // Concrete HTTP request structure
	Response  *HTTPResponse                 `json:"response,omitempty"`   // Concrete HTTP response structure
}

// HTTPRequest represents a recorded HTTP request.
type HTTPRequest struct {
	Method  string      `json:"method"`
	URL     string      `json:"url"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

// HTTPResponse represents a recorded HTTP response.
type HTTPResponse struct {
	Status     string      `json:"status"`
	StatusCode int         `json:"statusCode"`
	Headers    http.Header `json:"headers"`
	Body       []byte      `json:"body"`
}

// GeminiGenerateContentRequest captures the input parameters for a Gemini GenerateContent call.
type GeminiGenerateContentRequest struct {
	Model    string                     `json:"model"`
	Contents []*genai.Content           `json:"contents"`
	Config   *genai.GenerateContentConfig `json:"config,omitempty"`
}

// GeminiGenerateContentResponse captures the output from a Gemini GenerateContent call.
type GeminiGenerateContentResponse struct {
	Response *genai.GenerateContentResponse `json:"response"`
}
