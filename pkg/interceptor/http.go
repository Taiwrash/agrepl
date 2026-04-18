package interceptor

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"agrepl/pkg/core"
	"agrepl/pkg/storage"
)

// ANSI Color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// Mode defines the operation mode for the interceptor.
type Mode int

const (
	ModeRecord Mode = iota
	ModeReplay
	ModePassthrough
)

// HTTPInterceptor is an http.RoundTripper that can intercept and record/replay HTTP requests.
type HTTPInterceptor struct {
	Transport         http.RoundTripper
	Mode              Mode
	Storage           storage.Storage
	CurrentRun        *core.Run
	Fallback          bool // If true, call real network when no match is found
	IgnoreHeaders     []string
	IgnoreQueryParams []string
	usedSteps         map[int]bool // Keep track of which steps have been matched
	NetworkCallCount  int
	mu                sync.Mutex
}

// NewHTTPInterceptor creates a new HTTPInterceptor.
func NewHTTPInterceptor(mode Mode, s storage.Storage, currentRun *core.Run) *HTTPInterceptor {
	return &HTTPInterceptor{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Allow connecting to upstream even if certs are fishy (common in proxies)
		},
		Mode:              mode,
		Storage:           s,
		CurrentRun:        currentRun,
		usedSteps:         make(map[int]bool),
		IgnoreHeaders:     []string{"User-Agent", "Date", "Authorization", "X-Amz-Date"}, // Defaults
		IgnoreQueryParams: []string{"timestamp", "nonce", "api_key"},                    // Defaults
	}
}

// RoundTrip implements the http.RoundTripper interface.
func (i *HTTPInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	switch i.Mode {
	case ModeRecord:
		return i.recordRoundTrip(req)
	case ModeReplay:
		return i.replayRoundTrip(req)
	case ModePassthrough:
		fallthrough
	default:
		return i.Transport.RoundTrip(req)
	}
}

func (i *HTTPInterceptor) recordRoundTrip(req *http.Request) (*http.Response, error) {
	// Read the request body
	var reqBodyBytes []byte
	if req.Body != nil {
		reqBodyBytes, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
	}

	// Perform actual call
	resp, err := i.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Read response body
	var respBodyBytes []byte
	if resp.Body != nil {
		respBodyBytes, _ = ioutil.ReadAll(resp.Body)
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBodyBytes))
	}

	i.mu.Lock()
	i.NetworkCallCount++
	defer i.mu.Unlock()

	// Record step
	httpStep := core.Step{
		Type: core.StepTypeHTTP,
		Request: &core.HTTPRequest{
			Method:  req.Method,
			URL:     req.URL.String(),
			Headers: req.Header,
			Body:    reqBodyBytes,
		},
		Response: &core.HTTPResponse{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       respBodyBytes,
		},
	}
	i.CurrentRun.Steps = append(i.CurrentRun.Steps, httpStep)
	if i.Storage != nil {
		i.Storage.AppendStep(i.CurrentRun.RunID, httpStep)
	}
	fmt.Printf("%s[RECORD]%s Captured HTTP: %s %s\n", colorCyan, colorReset, req.Method, req.URL.String())

	return resp, nil
}

func (i *HTTPInterceptor) replayRoundTrip(req *http.Request) (*http.Response, error) {
	i.mu.Lock()

	// Read incoming request body
	var incomingReqBodyBytes []byte
	if req.Body != nil {
		incomingReqBodyBytes, _ = ioutil.ReadAll(req.Body)
		req.Body = ioutil.NopCloser(bytes.NewBuffer(incomingReqBodyBytes))
	}

	normalizedIncomingURL := i.normalizeURL(req.URL.String())

	// Search for first un-used match
	for idx, step := range i.CurrentRun.Steps {
		if i.usedSteps[idx] || step.Type != core.StepTypeHTTP {
			continue
		}

		recordedReq := step.Request
		if recordedReq == nil {
			continue
		}

		// Matcher logic
		methodMatch := recordedReq.Method == req.Method
		normalizedRecordedURL := i.normalizeURL(recordedReq.URL)
		urlMatch := normalizedRecordedURL == normalizedIncomingURL
		bodyMatch := i.compareBodies(recordedReq.Body, incomingReqBodyBytes)

		if methodMatch && urlMatch && bodyMatch {
			i.usedSteps[idx] = true
			recordedResp := step.Response
			fmt.Printf("%s[REPLAY] Matched request: %s %s%s\n", colorGreen, req.Method, req.URL.String(), colorReset)
			fmt.Printf("         Returning recorded response (Run: %s, Step: %d)\n", i.CurrentRun.RunID, idx)

			i.mu.Unlock()
			return &http.Response{
				StatusCode: recordedResp.StatusCode,
				Status:     recordedResp.Status,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     recordedResp.Headers,
				Body:       ioutil.NopCloser(bytes.NewBuffer(recordedResp.Body)),
				Request:    req,
			}, nil
		}
	}

	if i.Fallback {
		fmt.Printf("%s[REPLAY] No match found for %s %s. Falling back to real network.%s\n", colorYellow, req.Method, req.URL.String(), colorReset)
		i.NetworkCallCount++
		i.mu.Unlock()
		return i.Transport.RoundTrip(req)
	}

	fmt.Printf("%s[REPLAY] ERROR: No matching recorded step found for %s %s%s\n", colorRed, req.Method, req.URL.String(), colorReset)
	i.mu.Unlock()
	return nil, fmt.Errorf("replay error: no matching recorded step found for request %s %s", req.Method, req.URL.String())
}

func (i *HTTPInterceptor) normalizeURL(uStr string) string {
	u, err := url.Parse(uStr)
	if err != nil {
		return uStr
	}

	q := u.Query()
	for _, param := range i.IgnoreQueryParams {
		q.Del(param)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (i *HTTPInterceptor) compareBodies(recordedBody, incomingBody []byte) bool {
	if bytes.Equal(recordedBody, incomingBody) {
		return true
	}

	// Try JSON comparison
	var recordedJSON, incomingJSON interface{}
	err1 := json.Unmarshal(recordedBody, &recordedJSON)
	err2 := json.Unmarshal(incomingBody, &incomingJSON)

	if err1 == nil && err2 == nil {
		return reflect.DeepEqual(recordedJSON, incomingJSON)
	}

	return false
}
