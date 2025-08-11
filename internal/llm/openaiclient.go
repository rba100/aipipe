package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ModelType represents the type of model to use
type ModelType int

const (
	ModelTypeFast ModelType = iota
	ModelTypeDefault
	ModelTypeReasoning
	ModelTypeLocal
)

// Config holds the configuration for the LLM client
type Config struct {
	// API configuration
	APIEndpoint string
	APIToken    string

	// Model configuration
	DefaultModel   string
	FastModel      string
	ReasoningModel string
	LocalModel     string
	LocalBaseUrl   string

	// Common configuration
	IsCodeBlock bool
	IsStream    bool
	ModelType   ModelType
}

// LLMClient is the interface for interacting with LLM providers
type LLMClient interface {
	CreateCompletion(prompt string) (string, error)
	CreateCompletionStream(prompt string) <-chan string
}

// OpenAIClient implements the LLMClient interface for OpenAI/Groq
type OpenAIClient struct {
	config     *Config
	httpClient *http.Client
	baseURL    *url.URL
	apiKey     string
}

// NewClient creates a new LLM client
func NewClient(config *Config) (LLMClient, error) {
	if config.APIToken == "" {
		return nil, fmt.Errorf("API token is required")
	}

	var baseURL *url.URL
	var err error

	// Override base URL if provided
	if config.APIEndpoint != "" {
		baseURL, err = url.Parse(config.APIEndpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid API endpoint URL: %v", err)
		}
	} else {
		// Default to OpenAI endpoint if not specified
		baseURL, _ = url.Parse("https://api.openai.com/v1")
	}

	return &OpenAIClient{
		config:     config,
		httpClient: &http.Client{},
		baseURL:    baseURL,
		apiKey:     config.APIToken,
	}, nil
}

// GetModel returns the appropriate model based on the config
func (c *OpenAIClient) GetModel() string {
	switch c.config.ModelType {
	case ModelTypeFast:
		return c.config.FastModel
	case ModelTypeReasoning:
		return c.config.ReasoningModel
	case ModelTypeLocal:
		return c.config.LocalModel
	default:
		return c.config.DefaultModel
	}
}

// GetSystemPrompt returns the system prompt based on whether code block extraction is enabled
func GetSystemPrompt(isCodeBlock bool) string {
	if isCodeBlock {
		return "You are a helpful assistant. If the user has asked for something written, put it in a single code block (```type\\n...\\n```), otherwise just provide the answer."
	}
	return "You are a helpful assistant."
}

// CreateCompletion sends a prompt to the API and returns the completion
func (c *OpenAIClient) CreateCompletion(prompt string) (string, error) {
	model := c.GetModel()

	// Prepare the request body
	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": GetSystemPrompt(c.config.IsCodeBlock),
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Create the HTTP request
	endpoint := c.baseURL.String()
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	req, err := http.NewRequest("POST", endpoint+"chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "n/a" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	// Extract the completion text
	choices, ok := responseBody["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("invalid response format: missing choices")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format: invalid choice")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format: missing message")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid response format: missing content")
	}

	return content, nil
}

// CreateCompletionStream sends a prompt to the API and returns a stream of completions
func (c *OpenAIClient) CreateCompletionStream(prompt string) <-chan string {
	resultChan := make(chan string)
	errorChan := make(chan error, 1) // Buffer of 1 to avoid blocking

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		model := c.GetModel()

		// Prepare the request body
		requestBody := map[string]interface{}{
			"model": model,
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": GetSystemPrompt(c.config.IsCodeBlock),
				},
				{
					"role":    "user",
					"content": prompt,
				},
			},
			"stream": true,
		}

		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			errorChan <- fmt.Errorf("error marshaling request: %v", err)
			return
		}

		// Create the HTTP request
		endpoint := c.baseURL.String()
		if !strings.HasSuffix(endpoint, "/") {
			endpoint += "/"
		}
		req, err := http.NewRequest("POST", endpoint+"chat/completions", bytes.NewBuffer(jsonBody))
		if err != nil {
			errorChan <- fmt.Errorf("error creating request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		if c.apiKey != "n/a" {
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		}

		// Send the request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("error sending request: %v", err)
			return
		}
		defer resp.Body.Close()

		// Check for errors
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(bodyBytes))
			return
		}

		// Process the streaming response
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errorChan <- fmt.Errorf("error reading stream: %v", err)
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var streamResponse map[string]interface{}
			if err := json.Unmarshal([]byte(data), &streamResponse); err != nil {
				errorChan <- fmt.Errorf("error parsing stream data: %v", err)
				continue
			}

			choices, ok := streamResponse["choices"].([]interface{})
			if !ok || len(choices) == 0 {
				continue
			}

			choice, ok := choices[0].(map[string]interface{})
			if !ok {
				continue
			}

			delta, ok := choice["delta"].(map[string]interface{})
			if !ok {
				continue
			}

			content, ok := delta["content"].(string)
			if !ok || content == "" {
				continue
			}

			resultChan <- content
		}
	}()

	// Monitor the error channel and log errors
	go func() {
		for err := range errorChan {
			// Log the error to stderr
			fmt.Fprintf(os.Stderr, "Error in completion stream: %v\n", err)
		}
	}()

	return resultChan
}
