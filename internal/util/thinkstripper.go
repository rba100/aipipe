package util

import (
	"regexp"
	"strings"
)

// ThinkTagState represents the current state of the think tag processor
type ThinkTagState int

const (
	Searching ThinkTagState = iota
	Thinking
	Emitting
)

func StripThinkTags(input string) string {
	// Use (?s) to make dot match newlines, and use non-greedy matching with .*?
	thinkRegex := regexp.MustCompile(`(?s)^\s*<think>.*?</think>\s*`)
	result := thinkRegex.ReplaceAllString(input, "")
	return result
}

// StripThinkTagsStream processes a stream of text and strips think tags
func StripThinkTagsStream(inputStream <-chan string) <-chan string {
	startingRegex := regexp.MustCompile(`^[\s]*<think>`)
	outputStream := make(chan string)
	firstEmit := true

	go func() {
		defer close(outputStream)

		state := Searching
		var buffer strings.Builder

		for chunk := range inputStream {

			if state == Emitting {
				if firstEmit {
					chunk = strings.TrimLeft(chunk, " \n\t")
					if len(chunk) > 0 {
						outputStream <- chunk
						firstEmit = false
					}
				} else {
					outputStream <- chunk
				}
				continue
			}

			buffer.WriteString(chunk)
			currentBuffer := buffer.String()

			switch state {
			case Searching:
				if buffer.Len() < 10 {
					continue
				}
				if startingRegex.MatchString(currentBuffer) {
					state = Thinking
					closingIndex := strings.Index(currentBuffer, "</think>")
					if closingIndex != -1 {
						state = Emitting
						remainingText := currentBuffer[closingIndex+len("</think>"):]
						remainingText = strings.TrimLeft(remainingText, " \n\t")
						if len(remainingText) > 0 {
							outputStream <- remainingText
							firstEmit = false
						}
					}
					continue
				}
				state = Emitting
				outputStream <- currentBuffer
				firstEmit = false
				continue

			case Thinking:
				closingIndex := strings.Index(currentBuffer, "</think>")
				if closingIndex != -1 {
					state = Emitting
					remainingText := currentBuffer[closingIndex+len("</think>"):]
					remainingText = strings.TrimLeft(remainingText, " \n\t")
					if len(remainingText) > 0 {
						outputStream <- remainingText
						firstEmit = false
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
