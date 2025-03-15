package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAPIConfig(t *testing.T) {
	// Save original environment variables
	originalAipipeKey := os.Getenv("AIPIPE_API_KEY")
	originalGroqKey := os.Getenv("GROQ_API_KEY")
	originalOpenAIKey := os.Getenv("OPENAI_API_KEY")
	originalAipipeEndpoint := os.Getenv("AIPIPE_ENDPOINT")

	// Restore environment variables after test
	defer func() {
		os.Setenv("AIPIPE_API_KEY", originalAipipeKey)
		os.Setenv("GROQ_API_KEY", originalGroqKey)
		os.Setenv("OPENAI_API_KEY", originalOpenAIKey)
		os.Setenv("AIPIPE_ENDPOINT", originalAipipeEndpoint)
	}()

	// Test cases
	tests := []struct {
		name           string
		envVars        map[string]string
		expectError    bool
		expectedConfig *APIConfig
	}{
		{
			name: "AIPIPE_API_KEY set",
			envVars: map[string]string{
				"AIPIPE_API_KEY":  "test-aipipe-key",
				"AIPIPE_ENDPOINT": "https://test-aipipe-endpoint.com",
				"GROQ_API_KEY":    "",
				"OPENAI_API_KEY":  "",
			},
			expectError: false,
			expectedConfig: &APIConfig{
				APIToken:       "test-aipipe-key",
				APIEndpoint:    "https://test-aipipe-endpoint.com",
				DefaultModel:   "llama-3.3-70b-versatile",
				FastModel:      "llama-3.1-8b-instant",
				ReasoningModel: "qwen-qwq-32b",
			},
		},
		{
			name: "GROQ_API_KEY set",
			envVars: map[string]string{
				"AIPIPE_API_KEY":  "",
				"AIPIPE_ENDPOINT": "",
				"GROQ_API_KEY":    "test-groq-key",
				"OPENAI_API_KEY":  "",
			},
			expectError: false,
			expectedConfig: &APIConfig{
				APIToken:       "test-groq-key",
				APIEndpoint:    "https://api.groq.com/openai/v1",
				DefaultModel:   "llama-3.3-70b-versatile",
				FastModel:      "llama-3.1-8b-instant",
				ReasoningModel: "qwen-qwq-32b",
			},
		},
		{
			name: "OPENAI_API_KEY set",
			envVars: map[string]string{
				"AIPIPE_API_KEY":  "",
				"AIPIPE_ENDPOINT": "",
				"GROQ_API_KEY":    "",
				"OPENAI_API_KEY":  "test-openai-key",
			},
			expectError: false,
			expectedConfig: &APIConfig{
				APIToken:       "test-openai-key",
				APIEndpoint:    "https://api.openai.com/v1",
				DefaultModel:   "gpt-4o",
				FastModel:      "gpt-4o-mini",
				ReasoningModel: "o3-mini",
			},
		},
		{
			name: "No API keys set",
			envVars: map[string]string{
				"AIPIPE_API_KEY":  "",
				"AIPIPE_ENDPOINT": "",
				"GROQ_API_KEY":    "",
				"OPENAI_API_KEY":  "",
			},
			expectError: true,
		},
		{
			name: "AIPIPE_API_KEY set but no endpoint",
			envVars: map[string]string{
				"AIPIPE_API_KEY":  "test-aipipe-key",
				"AIPIPE_ENDPOINT": "",
				"GROQ_API_KEY":    "",
				"OPENAI_API_KEY":  "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables for this test
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Call GetAPIConfig
			config, err := GetAPIConfig()

			// Check error
			if tt.expectError {
				if err == nil {
					t.Errorf("GetAPIConfig() error = nil, expected an error")
				}
				return
			}

			if err != nil {
				t.Errorf("GetAPIConfig() error = %v, expected no error", err)
				return
			}

			// Check config values
			if config.APIToken != tt.expectedConfig.APIToken {
				t.Errorf("APIToken = %v, want %v", config.APIToken, tt.expectedConfig.APIToken)
			}
			if config.APIEndpoint != tt.expectedConfig.APIEndpoint {
				t.Errorf("APIEndpoint = %v, want %v", config.APIEndpoint, tt.expectedConfig.APIEndpoint)
			}
			if config.DefaultModel != tt.expectedConfig.DefaultModel {
				t.Errorf("DefaultModel = %v, want %v", config.DefaultModel, tt.expectedConfig.DefaultModel)
			}
			if config.FastModel != tt.expectedConfig.FastModel {
				t.Errorf("FastModel = %v, want %v", config.FastModel, tt.expectedConfig.FastModel)
			}
			if config.ReasoningModel != tt.expectedConfig.ReasoningModel {
				t.Errorf("ReasoningModel = %v, want %v", config.ReasoningModel, tt.expectedConfig.ReasoningModel)
			}
		})
	}
}

func TestLoadUserConfig(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "aipipe-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Save original home directory environment variables
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	// Set home directory environment variables to our temp directory for the test
	os.Setenv("HOME", tempDir)
	os.Setenv("USERPROFILE", tempDir)

	// Restore home directory environment variables after test
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("USERPROFILE", originalUserProfile)
	}()

	// Create .aipipe directory
	aipipeDir := filepath.Join(tempDir, ".aipipe")
	if err := os.MkdirAll(aipipeDir, 0755); err != nil {
		t.Fatalf("Failed to create .aipipe dir: %v", err)
	}

	// Test cases
	tests := []struct {
		name           string
		configContent  string
		initialConfig  *APIConfig
		expectedConfig *APIConfig
		expectError    bool
	}{
		{
			name: "Valid config file",
			configContent: `
endpoint: https://custom-endpoint.com
apiKey: custom-api-key
defaultModel: custom-default-model
fastModel: custom-fast-model
reasoningModel: custom-reasoning-model
`,
			initialConfig: &APIConfig{
				APIToken:       "initial-token",
				APIEndpoint:    "initial-endpoint",
				DefaultModel:   "initial-default-model",
				FastModel:      "initial-fast-model",
				ReasoningModel: "initial-reasoning-model",
			},
			expectedConfig: &APIConfig{
				APIToken:       "custom-api-key",
				APIEndpoint:    "https://custom-endpoint.com",
				DefaultModel:   "custom-default-model",
				FastModel:      "custom-fast-model",
				ReasoningModel: "custom-reasoning-model",
			},
			expectError: false,
		},
		{
			name: "Case insensitive keys",
			configContent: `
ENDPOINT: https://custom-endpoint.com
ApiKey: custom-api-key
DefaultModel: custom-default-model
fastmodel: custom-fast-model
reasoningMODEL: custom-reasoning-model
`,
			initialConfig: &APIConfig{
				APIToken:       "initial-token",
				APIEndpoint:    "initial-endpoint",
				DefaultModel:   "initial-default-model",
				FastModel:      "initial-fast-model",
				ReasoningModel: "initial-reasoning-model",
			},
			expectedConfig: &APIConfig{
				APIToken:       "custom-api-key",
				APIEndpoint:    "https://custom-endpoint.com",
				DefaultModel:   "custom-default-model",
				FastModel:      "custom-fast-model",
				ReasoningModel: "custom-reasoning-model",
			},
			expectError: false,
		},
		{
			name: "Partial config file",
			configContent: `
endpoint: https://custom-endpoint.com
`,
			initialConfig: &APIConfig{
				APIToken:       "initial-token",
				APIEndpoint:    "initial-endpoint",
				DefaultModel:   "initial-default-model",
				FastModel:      "initial-fast-model",
				ReasoningModel: "initial-reasoning-model",
			},
			expectedConfig: &APIConfig{
				APIToken:       "initial-token",
				APIEndpoint:    "https://custom-endpoint.com",
				DefaultModel:   "initial-default-model",
				FastModel:      "initial-fast-model",
				ReasoningModel: "initial-reasoning-model",
			},
			expectError: false,
		},
		{
			name:          "Invalid YAML",
			configContent: `invalid: yaml: :`,
			initialConfig: &APIConfig{
				APIToken:       "initial-token",
				APIEndpoint:    "initial-endpoint",
				DefaultModel:   "initial-default-model",
				FastModel:      "initial-fast-model",
				ReasoningModel: "initial-reasoning-model",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config file
			configPath := filepath.Join(aipipeDir, "config.yaml")
			if err := os.WriteFile(configPath, []byte(tt.configContent), 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			// Create a copy of the initial config
			config := &APIConfig{
				APIToken:       tt.initialConfig.APIToken,
				APIEndpoint:    tt.initialConfig.APIEndpoint,
				DefaultModel:   tt.initialConfig.DefaultModel,
				FastModel:      tt.initialConfig.FastModel,
				ReasoningModel: tt.initialConfig.ReasoningModel,
			}

			// Call LoadUserConfig
			err := LoadUserConfig(config)

			// Check error
			if tt.expectError {
				if err == nil {
					t.Errorf("LoadUserConfig() error = nil, expected an error")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadUserConfig() error = %v, expected no error", err)
				return
			}

			// Check config values
			if config.APIToken != tt.expectedConfig.APIToken {
				t.Errorf("APIToken = %v, want %v", config.APIToken, tt.expectedConfig.APIToken)
			}
			if config.APIEndpoint != tt.expectedConfig.APIEndpoint {
				t.Errorf("APIEndpoint = %v, want %v", config.APIEndpoint, tt.expectedConfig.APIEndpoint)
			}
			if config.DefaultModel != tt.expectedConfig.DefaultModel {
				t.Errorf("DefaultModel = %v, want %v", config.DefaultModel, tt.expectedConfig.DefaultModel)
			}
			if config.FastModel != tt.expectedConfig.FastModel {
				t.Errorf("FastModel = %v, want %v", config.FastModel, tt.expectedConfig.FastModel)
			}
			if config.ReasoningModel != tt.expectedConfig.ReasoningModel {
				t.Errorf("ReasoningModel = %v, want %v", config.ReasoningModel, tt.expectedConfig.ReasoningModel)
			}
		})
	}
}
