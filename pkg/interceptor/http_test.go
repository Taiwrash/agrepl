package interceptor

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"agrepl/pkg/core"
)

func TestHTTPInterceptor_ReplayMatching(t *testing.T) {
	run := &core.Run{
		RunID: "test-run",
		Steps: []core.Step{
			{
				Type: core.StepTypeHTTP,
				Request: &core.HTTPRequest{
					Method: "GET",
					URL:    "https://example.com/api/v1/resource",
					Body:   "",
				},
				Response: &core.HTTPResponse{
					Status:     "200 OK",
					StatusCode: 200,
					Body:       "{\"id\": 1}",
				},
			},
			{
				Type: core.StepTypeHTTP,
				Request: &core.HTTPRequest{
					Method: "POST",
					URL:    "https://example.com/api/v1/resource",
					Body:   "{\"name\": \"test\"}",
				},
				Response: &core.HTTPResponse{
					Status:     "201 Created",
					StatusCode: 201,
					Body:       "{\"id\": 2}",
				},
			},
		},
	}

	interceptor := NewHTTPInterceptor(ModeReplay, nil, run)

	t.Run("Match GET", func(t *testing.T) {
		req := &http.Request{
			Method: "GET",
			URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/api/v1/resource"},
		}
		resp, err := interceptor.RoundTrip(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != "{\"id\": 1}" {
			t.Errorf("Expected body {\"id\": 1}, got %s", string(body))
		}
	})

	t.Run("Match POST", func(t *testing.T) {
		req := &http.Request{
			Method: "POST",
			URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/api/v1/resource"},
			Body:   ioutil.NopCloser(bytes.NewBufferString("{\"name\": \"test\"}")),
		}
		resp, err := interceptor.RoundTrip(req)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if resp.StatusCode != 201 {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != "{\"id\": 2}" {
			t.Errorf("Expected body {\"id\": 2}, got %s", string(body))
		}
	})

	t.Run("Match with Ignore Query Params", func(t *testing.T) {
		runWithQuery := &core.Run{
			Steps: []core.Step{
				{
					Type: core.StepTypeHTTP,
					Request: &core.HTTPRequest{
						Method: "GET",
						URL:    "https://example.com/api?api_key=secret&timestamp=123",
					},
					Response: &core.HTTPResponse{StatusCode: 200},
				},
			},
		}
		interceptor := NewHTTPInterceptor(ModeReplay, nil, runWithQuery)
		interceptor.IgnoreQueryParams = []string{"timestamp", "api_key"}

		req := &http.Request{
			Method: "GET",
			URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/api", RawQuery: "api_key=different&timestamp=456"},
		}
		resp, err := interceptor.RoundTrip(req)
		if err != nil {
			t.Fatalf("Expected match with ignored query params, got error: %v", err)
		}
		if resp.StatusCode != 200 {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("No Match", func(t *testing.T) {
		req := &http.Request{
			Method: "GET",
			URL:    &url.URL{Scheme: "https", Host: "example.com", Path: "/api/v1/unknown"},
		}
		_, err := interceptor.RoundTrip(req)
		if err == nil {
			t.Fatal("Expected error for no match, got nil")
		}
	})
}
