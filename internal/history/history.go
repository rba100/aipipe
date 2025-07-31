package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AIPipeDir is the name of the directory where we store config and history.
const AIPipeDir = ".aipipe"

// LastConversationFile is the name of the file holding the last conversation.
const LastConversationFile = "last-conversation.json"

// HistoryDir is the name of the directory where old conversations are archived.
const HistoryDir = "history"

// Paths holds the important paths for history management.
type Paths struct {
	BaseDir      string
	HistoryDir   string
	LastConvFile string
}

// GetPaths returns the key paths for managing conversation history.
func GetPaths() (*Paths, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	baseDir := filepath.Join(homeDir, AIPipeDir)
	historyDir := filepath.Join(baseDir, HistoryDir)
	lastConvFile := filepath.Join(baseDir, LastConversationFile)

	return &Paths{
		BaseDir:      baseDir,
		HistoryDir:   historyDir,
		LastConvFile: lastConvFile,
	}, nil
}

// ReadConversation reads a conversation from a JSON file.
func ReadConversation(path string) (*Conversation, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the file doesn't exist, return an empty conversation.
		return &Conversation{Messages: []Message{}}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read conversation file: %w", err)
	}

	var conv Conversation
	if err := json.Unmarshal(data, &conv); err != nil {
		return nil, fmt.Errorf("failed to parse conversation file: %w", err)
	}

	return &conv, nil
}

// WriteConversation writes a conversation to a JSON file.
func WriteConversation(path string, conversation *Conversation) error {
	data, err := json.MarshalIndent(conversation, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal conversation: %w", err)
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return os.WriteFile(path, data, 0644)
}

// ArchiveLastConversation moves the last conversation to the history directory.
func ArchiveLastConversation() error {
	paths, err := GetPaths()
	if err != nil {
		return err
	}

	if _, err := os.Stat(paths.LastConvFile); os.IsNotExist(err) {
		// No last conversation to archive.
		return nil
	}

	// Ensure history directory exists.
	if _, err := os.Stat(paths.HistoryDir); os.IsNotExist(err) {
		if err := os.MkdirAll(paths.HistoryDir, 0755); err != nil {
			return fmt.Errorf("failed to create history directory: %w", err)
		}
	}

	// Move the file.
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	newPath := filepath.Join(paths.HistoryDir, fmt.Sprintf("%s_%s", timestamp, LastConversationFile))

	if err := os.Rename(paths.LastConvFile, newPath); err != nil {
		return fmt.Errorf("failed to move last conversation file: %w", err)
	}

	return nil
}
