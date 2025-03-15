package util

import (
	"regexp"
	"strings"
)

// CodeBlockResult represents the result of extracting a code block
type CodeBlockResult struct {
	Text string
	Type string
}

// ExtractCodeBlock extracts a code block from a string
func ExtractCodeBlock(input string) CodeBlockResult {
	// Use a regex pattern that can handle empty code blocks and capture the language type
	re := regexp.MustCompile("```([a-zA-Z0-9.]*)(?:\n)?([\\s\\S]*?)(?:\n```|```)")
	matches := re.FindStringSubmatch(input)
	if len(matches) > 2 {
		return CodeBlockResult{
			Text: matches[2],
			Type: matches[1],
		}
	}
	return CodeBlockResult{
		Text: input,
		Type: "",
	}
}

// CodeBlockState represents the state of code block extraction
type CodeBlockState int

const (
	SearchingOpening CodeBlockState = iota
	Open
	Closed
)

// ExtractCodeBlockStream extracts code blocks from a stream
func ExtractCodeBlockStream(inputStream <-chan string) <-chan CodeBlockResult {
	outputStream := make(chan CodeBlockResult)
	openingRe := regexp.MustCompile("```([a-zA-Z0-9.]*)(?:\n)")
	potentialClosingRe := regexp.MustCompile("\n`{0,2}$")

	go func() {
		defer close(outputStream)

		buffer := strings.Builder{}
		state := SearchingOpening
		var blockType string = ""

		for part := range inputStream {
			if state == Closed {
				break
			}

			buffer.WriteString(part)
			bufStr := buffer.String()

			if state == SearchingOpening {
				// Look for opening marker with optional language type
				match := openingRe.FindStringSubmatchIndex(bufStr)

				if len(match) > 0 {
					// Extract the language type if present
					if match[2] != -1 && match[3] != -1 {
						blockType = bufStr[match[2]:match[3]]
					}

					// Move to the content after the opening marker
					remainingContent := bufStr[match[1]:]
					buffer.Reset()
					buffer.WriteString(remainingContent)
					state = Open
					continue
				}
			}

			if state == Open {
				// Check for potential closing marker at the end
				if potentialClosingRe.MatchString(bufStr) {
					continue
				}

				// Check for actual closing marker
				closePos := strings.Index(bufStr, "\n```")
				if closePos >= 0 || strings.HasPrefix(bufStr, "```") {
					output := bufStr[:closePos]
					state = Closed
					buffer.Reset()
					outputStream <- CodeBlockResult{
						Text: output,
						Type: blockType,
					}
					break
				}

				// If we're still processing and have content, return it and clear buffer
				output := bufStr
				buffer.Reset()
				outputStream <- CodeBlockResult{
					Text: output,
					Type: blockType,
				}
			}
		}

		// If we never closed the code block but have content, return what we have
		if state != Closed && buffer.Len() > 0 {
			remainingContent := buffer.String()

			if strings.HasPrefix(remainingContent, "```") {
				return
			}

			outputStream <- CodeBlockResult{
				Text: remainingContent,
				Type: blockType,
			}
		}
	}()

	return outputStream
}
