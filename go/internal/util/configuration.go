package util

import (
	"fmt"
	"os"
)

// APIConfig holds the configuration for the API client
type APIConfig struct {
	APIEndpoint string
	APIToken    string
	IsAipipe    bool
	IsGroq      bool
	IsOpenAI    bool
}

// GetAPIConfig retrieves API configuration from environment variables
func GetAPIConfig() (*APIConfig, error) {
	config := &APIConfig{}

	// Check for API keys in order of preference
	config.APIToken = os.Getenv("AIPIPE_API_KEY")
	if config.APIToken != "" {
		config.IsAipipe = true
	}

	if config.APIToken == "" {
		config.APIToken = os.Getenv("GROQ_API_KEY")
		config.IsGroq = true
	}

	if config.APIToken == "" {
		config.APIToken = os.Getenv("OPENAI_API_KEY")
		config.IsOpenAI = true
	}

	if config.APIToken == "" {
		return nil, fmt.Errorf("AIPIPE_API_KEY or GROQ_API_KEY or OPENAI_API_KEY environment variable is not set")
	}

	// Set API endpoint based on the service type
	config.APIEndpoint = os.Getenv("AIPIPE_ENDPOINT")
	if config.IsAipipe && config.APIEndpoint == "" {
		return nil, fmt.Errorf("AIPIPE_ENDPOINT environment variable is not set")
	}

	if config.IsOpenAI {
		config.APIEndpoint = "https://api.openai.com/v1"
	}

	if config.IsGroq {
		config.APIEndpoint = "https://api.groq.com/openai/v1"
	}

	return config, nil
}
