package llm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestGetModel tests the GetModel function
func TestGetModel(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name: "Default model",
			config: &Config{
				ModelType:      ModelTypeDefault,
				DefaultModel:   "default-model",
				FastModel:      "fast-model",
				ReasoningModel: "reasoning-model",
			},
			expected: "default-model",
		},
		{
			name: "Fast model",
			config: &Config{
				ModelType:      ModelTypeFast,
				DefaultModel:   "default-model",
				FastModel:      "fast-model",
				ReasoningModel: "reasoning-model",
			},
			expected: "fast-model",
		},
		{
			name: "Reasoning model",
			config: &Config{
				ModelType:      ModelTypeReasoning,
				DefaultModel:   "default-model",
				FastModel:      "fast-model",
				ReasoningModel: "reasoning-model",
			},
			expected: "reasoning-model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &OpenAIClient{
				config: tt.config,
			}
			if got := client.GetModel(); got != tt.expected {
				t.Errorf("GetModel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestGetSystemPrompt tests the GetSystemPrompt function
func TestGetSystemPrompt(t *testing.T) {
	tests := []struct {
		name        string
		isCodeBlock bool
		expected    string
	}{
		{
			name:        "With code block",
			isCodeBlock: true,
			expected:    "You are a helpful assistant. If the user has asked for something written, put it in a single code block (```type\\n...\\n```), otherwise just provide the answer.",
		},
		{
			name:        "Without code block",
			isCodeBlock: false,
			expected:    "You are a helpful assistant.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSystemPrompt(tt.isCodeBlock); got != tt.expected {
				t.Errorf("GetSystemPrompt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "Valid config",
			config: &Config{
				APIToken: "test-token",
			},
			expectError: false,
		},
		{
			name: "Missing API token",
			config: &Config{
				APIToken: "",
			},
			expectError: true,
		},
		{
			name: "Custom endpoint",
			config: &Config{
				APIToken:    "test-token",
				APIEndpoint: "https://custom-api.example.com",
			},
			expectError: false,
		},
		{
			name: "Invalid endpoint",
			config: &Config{
				APIToken:    "test-token",
				APIEndpoint: "://invalid-url",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("NewClient() error = nil, expected an error")
				}
			} else {
				if err != nil {
					t.Errorf("NewClient() error = %v, expected no error", err)
				}
				if client == nil {
					t.Errorf("NewClient() client is nil, expected a client")
				}
			}
		})
	}
}

// mockHTTPClient is a mock HTTP client for testing
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestCreateCompletion tests the CreateCompletion function
func TestCreateCompletion(t *testing.T) {
	// Test successful completion
	t.Run("Successful completion", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request method
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Check authorization header
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
			}

			// Check content type
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type header 'application/json', got %s", r.Header.Get("Content-Type"))
			}

			// Read request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Error reading request body: %v", err)
			}

			// Parse request body
			var requestBody map[string]interface{}
			if err := json.Unmarshal(body, &requestBody); err != nil {
				t.Errorf("Error parsing request body: %v", err)
			}

			// Check model
			if requestBody["model"] != "test-model" {
				t.Errorf("Expected model 'test-model', got %v", requestBody["model"])
			}

			// Write response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"choices": [
					{
						"message": {
							"content": "Test response"
						}
					}
				]
			}`))
		}))
		defer server.Close()

		// Create client
		baseURL, _ := url.Parse(server.URL)
		client := &OpenAIClient{
			config: &Config{
				DefaultModel: "test-model",
				ModelType:    ModelTypeDefault,
			},
			httpClient: server.Client(),
			baseURL:    baseURL,
			apiKey:     "test-token",
		}

		// Call CreateCompletion
		response, err := client.CreateCompletion("Test prompt")
		if err != nil {
			t.Errorf("CreateCompletion() error = %v, expected no error", err)
		}
		if response != "Test response" {
			t.Errorf("CreateCompletion() = %v, want %v", response, "Test response")
		}
	})

	// Test error response
	t.Run("Error response", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Test error"}`))
		}))
		defer server.Close()

		// Create client
		baseURL, _ := url.Parse(server.URL)
		client := &OpenAIClient{
			config: &Config{
				DefaultModel: "test-model",
				ModelType:    ModelTypeDefault,
			},
			httpClient: server.Client(),
			baseURL:    baseURL,
			apiKey:     "test-token",
		}

		// Call CreateCompletion
		_, err := client.CreateCompletion("Test prompt")
		if err == nil {
			t.Errorf("CreateCompletion() error = nil, expected an error")
		}
	})

	// Test invalid response format
	t.Run("Invalid response format", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`)) // Empty response
		}))
		defer server.Close()

		// Create client
		baseURL, _ := url.Parse(server.URL)
		client := &OpenAIClient{
			config: &Config{
				DefaultModel: "test-model",
				ModelType:    ModelTypeDefault,
			},
			httpClient: server.Client(),
			baseURL:    baseURL,
			apiKey:     "test-token",
		}

		// Call CreateCompletion
		_, err := client.CreateCompletion("Test prompt")
		if err == nil {
			t.Errorf("CreateCompletion() error = nil, expected an error")
		}
	})
}

// TestCreateCompletionStream tests the CreateCompletionStream function
func TestCreateCompletionStream(t *testing.T) {
	// Test successful stream
	t.Run("Successful stream", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check request method
			if r.Method != "POST" {
				t.Errorf("Expected POST request, got %s", r.Method)
			}

			// Check authorization header
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("Expected Authorization header 'Bearer test-token', got %s", r.Header.Get("Authorization"))
			}

			// Check content type
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type header 'application/json', got %s", r.Header.Get("Content-Type"))
			}

			// Read request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Error reading request body: %v", err)
			}

			// Parse request body
			var requestBody map[string]interface{}
			if err := json.Unmarshal(body, &requestBody); err != nil {
				t.Errorf("Error parsing request body: %v", err)
			}

			// Check model
			if requestBody["model"] != "test-model" {
				t.Errorf("Expected model 'test-model', got %v", requestBody["model"])
			}

			// Check stream flag
			if requestBody["stream"] != true {
				t.Errorf("Expected stream flag to be true, got %v", requestBody["stream"])
			}

			// Write response
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			// Write stream data
			flusher, ok := w.(http.Flusher)
			if !ok {
				t.Errorf("ResponseWriter does not implement http.Flusher")
				return
			}

			// Send multiple chunks
			w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Part 1\"}}]}\n\n"))
			flusher.Flush()

			w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Part 2\"}}]}\n\n"))
			flusher.Flush()

			w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"Part 3\"}}]}\n\n"))
			flusher.Flush()

			w.Write([]byte("data: [DONE]\n\n"))
			flusher.Flush()
		}))
		defer server.Close()

		// Create client
		baseURL, _ := url.Parse(server.URL)
		client := &OpenAIClient{
			config: &Config{
				DefaultModel: "test-model",
				ModelType:    ModelTypeDefault,
			},
			httpClient: server.Client(),
			baseURL:    baseURL,
			apiKey:     "test-token",
		}

		// Call CreateCompletionStream
		stream := client.CreateCompletionStream("Test prompt")

		// Collect stream results
		var results []string
		for part := range stream {
			results = append(results, part)
		}

		// Check results
		expected := []string{"Part 1", "Part 2", "Part 3"}
		if len(results) != len(expected) {
			t.Errorf("CreateCompletionStream() returned %d parts, expected %d", len(results), len(expected))
		}

		for i, result := range results {
			if result != expected[i] {
				t.Errorf("CreateCompletionStream() part %d = %v, want %v", i, result, expected[i])
			}
		}
	})

	// Test error response
	t.Run("Error response", func(t *testing.T) {
		// Create a test server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Test error"}`))
		}))
		defer server.Close()

		// Create client
		baseURL, _ := url.Parse(server.URL)
		client := &OpenAIClient{
			config: &Config{
				DefaultModel: "test-model",
				ModelType:    ModelTypeDefault,
			},
			httpClient: server.Client(),
			baseURL:    baseURL,
			apiKey:     "test-token",
		}

		// Call CreateCompletionStream
		stream := client.CreateCompletionStream("Test prompt")

		// Collect stream results
		var results []string
		for part := range stream {
			results = append(results, part)
		}

		// Check results - should be empty since there was an error
		if len(results) != 0 {
			t.Errorf("CreateCompletionStream() returned %d parts, expected 0", len(results))
		}
	})
}
