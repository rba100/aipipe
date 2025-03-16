package util

import "strings"

// ThinkTagState represents the current state of the think tag processor
type ThinkTagState int

const (
	Searching ThinkTagState = iota
	Thinking
	Emitting
)

// StripThinkTagsStream processes a stream of text and strips think tags
func StripThinkTagsStream(inputStream <-chan string) <-chan string {
	outputStream := make(chan string)

	go func() {
		defer close(outputStream)

		state := Searching
		var buffer strings.Builder

		for chunk := range inputStream {

			if state == Emitting {
				outputStream <- chunk
				continue
			}

			buffer.WriteString(chunk)
			currentBuffer := buffer.String()

			switch state {
			case Searching:
				if buffer.Len() < 7 {
					continue
				}
				if strings.HasPrefix(currentBuffer, "<think>") {
					state = Thinking
					closingIndex := strings.Index(currentBuffer, "</think>")
					if closingIndex != -1 {
						state = Emitting
						remainingText := currentBuffer[closingIndex+len("</think>"):]
						if len(remainingText) > 0 {
							outputStream <- remainingText
						}
					}
					continue
				}
				state = Emitting
				outputStream <- currentBuffer
				continue

			case Thinking:
				closingIndex := strings.Index(currentBuffer, "</think>")
				if closingIndex != -1 {
					state = Emitting
					remainingText := currentBuffer[closingIndex+len("</think>"):]
					if len(remainingText) > 0 {
						outputStream <- remainingText
					}
					continue
				}
			}
		}

		// End of stream handling
		if state != Emitting {
			// If we're still in thinking mode, emit with <think> prefix
			outputStream <- buffer.String()
		}
	}()

	return outputStream
}
