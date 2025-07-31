package history

// Message represents a single message in a conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Conversation represents a series of messages.
type Conversation struct {
	Messages []Message `json:"messages"`
}
