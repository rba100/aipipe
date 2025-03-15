package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ModelType represents the type of model to use
type ModelType int

const (
	ModelTypeFast ModelType = iota
	ModelTypeDefault
	ModelTypeReasoning
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
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

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

	go func() {
		defer close(resultChan)

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
			fmt.Fprintf(io.Discard, "error marshaling request: %v", err)
			return
		}

		// Create the HTTP request
		endpoint := c.baseURL.String()
		if !strings.HasSuffix(endpoint, "/") {
			endpoint += "/"
		}
		req, err := http.NewRequest("POST", endpoint+"chat/completions", bytes.NewBuffer(jsonBody))
		if err != nil {
			fmt.Fprintf(io.Discard, "error creating request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiKey)

		// Send the request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(io.Discard, "error sending request: %v", err)
			return
		}
		defer resp.Body.Close()

		// Check for errors
		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(io.Discard, "API error (status %d): %s", resp.StatusCode, string(bodyBytes))
			return
		}

		// Process the streaming response
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Fprintf(io.Discard, "error reading stream: %v", err)
				}
				break
			}

			line = strings.TrimSpace(line)
			if line == "" || line == "data: [DONE]" {
				continue
			}

			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := line[6:] // Remove "data: " prefix

			var streamResponse map[string]interface{}
			if err := json.Unmarshal([]byte(data), &streamResponse); err != nil {
				fmt.Fprintf(io.Discard, "error parsing stream data: %v", err)
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

	return resultChan
}

/*
IMPLEMENTATION NOTES:

To use the official OpenAI Go client with Groq, you would need to:

1. Import the required packages:
   ```go
   import (
       "context"
       "github.com/openai/openai-go"
       "github.com/openai/openai-go/option"
   )
   ```

2. Initialize the client with the appropriate options:
   ```go
   options := []option.ClientOption{
       option.WithAPIKey(config.APIToken),
   }

   if baseURL != nil {
       options = append(options, option.WithBaseURL(baseURL))
   }

   client := openai.NewClient(options...)
   ```

3. Use the client to create completions:
   ```go
   chatCompletion, err := client.Chat.Completions.New(
       context.Background(),
       openai.ChatCompletionNewParams{
           Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
               openai.SystemMessage(GetSystemPrompt(c.config.IsCodeBlock)),
               openai.UserMessage(prompt),
           }),
           Model: openai.F(model),
       },
   )
   ```

4. For streaming completions:
   ```go
   stream := client.Chat.Completions.NewStreaming(
       context.Background(),
       openai.ChatCompletionNewParams{
           Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
               openai.SystemMessage(GetSystemPrompt(c.config.IsCodeBlock)),
               openai.UserMessage(prompt),
           }),
           Model: openai.F(model),
       },
   )

   for stream.Next() {
       chunk := stream.Current()
       if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
           resultChan <- chunk.Choices[0].Delta.Content
       }
   }

   if err := stream.Err(); err != nil {
       // Handle error
   }
   ```
*/
