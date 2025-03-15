package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// APIConfig holds the configuration for the API client
type APIConfig struct {
	APIEndpoint    string
	APIToken       string
	DefaultModel   string
	FastModel      string
	ReasoningModel string
}

// UserConfig holds the user's configuration from YAML file
type UserConfig struct {
	Endpoint       string `yaml:"endpoint"`
	APIKey         string `yaml:"apiKey"`
	DefaultModel   string `yaml:"defaultModel"`
	FastModel      string `yaml:"fastModel"`
	ReasoningModel string `yaml:"reasoningModel"`
}

// LoadUserConfig loads configuration from ~/.aipipe/config.yaml if it exists
// and merges it with the existing APIConfig
func LoadUserConfig(config *APIConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".aipipe", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, just return without error
		return nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Create a map for case-insensitive parsing
	var configMap map[string]interface{}
	if err := yaml.Unmarshal(data, &configMap); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	// Convert keys to lowercase for case-insensitive matching
	normalizedMap := make(map[string]interface{})
	for k, v := range configMap {
		normalizedMap[strings.ToLower(k)] = v
	}

	// Extract values with case-insensitive keys
	if endpoint, ok := normalizedMap["endpoint"]; ok && endpoint != "" {
		if str, ok := endpoint.(string); ok {
			config.APIEndpoint = str
		}
	}

	if apiKey, ok := normalizedMap["apikey"]; ok && apiKey != "" {
		if str, ok := apiKey.(string); ok {
			config.APIToken = str
		}
	}

	if defaultModel, ok := normalizedMap["defaultmodel"]; ok && defaultModel != "" {
		if str, ok := defaultModel.(string); ok {
			config.DefaultModel = str
		}
	}

	if fastModel, ok := normalizedMap["fastmodel"]; ok && fastModel != "" {
		if str, ok := fastModel.(string); ok {
			config.FastModel = str
		}
	}

	if reasoningModel, ok := normalizedMap["reasoningmodel"]; ok && reasoningModel != "" {
		if str, ok := reasoningModel.(string); ok {
			config.ReasoningModel = str
		}
	}

	return nil
}

// GetAPIConfig retrieves API configuration from environment variables and config file
func GetAPIConfig() (*APIConfig, error) {
	config := &APIConfig{}

	isAipipe := false
	isGroq := false
	isOpenAI := false

	// Check for API keys in order of preference
	config.APIToken = os.Getenv("AIPIPE_API_KEY")
	if config.APIToken != "" {
		isAipipe = true
	}

	if config.APIToken == "" {
		config.APIToken = os.Getenv("GROQ_API_KEY")
		if config.APIToken != "" {
			isGroq = true
		}
	}

	if config.APIToken == "" {
		config.APIToken = os.Getenv("OPENAI_API_KEY")
		if config.APIToken != "" {
			isOpenAI = true
		}
	}

	config.DefaultModel = "llama-3.3-70b-versatile"
	config.FastModel = "llama-3.1-8b-instant"
	config.ReasoningModel = "qwen-2.5-32b"

	if isOpenAI {
		config.DefaultModel = "gpt-4o"
		config.FastModel = "gpt-4o-mini"
		config.ReasoningModel = "o3-mini"
	}

	// Try to load configuration from YAML file
	// This will override environment variables if values are present in the file
	if err := LoadUserConfig(config); err != nil {
		// Just log the error but continue with env vars
		fmt.Fprintf(os.Stderr, "Warning: Failed to load user config: %v\n", err)
	}

	// Final check if we have an API token
	if config.APIToken == "" {
		return nil, fmt.Errorf("AIPIPE_API_KEY or GROQ_API_KEY or OPENAI_API_KEY environment variable is not set and no API key found in config file")
	}

	// Set API endpoint based on the service type if not already set
	if config.APIEndpoint == "" {
		config.APIEndpoint = os.Getenv("AIPIPE_ENDPOINT")
		if isAipipe && config.APIEndpoint == "" {
			return nil, fmt.Errorf("AIPIPE_ENDPOINT environment variable is not set and no endpoint found in config file")
		}

		if isOpenAI {
			config.APIEndpoint = "https://api.openai.com/v1"
		}

		if isGroq {
			config.APIEndpoint = "https://api.groq.com/openai/v1"
		}
	}

	return config, nil
}
