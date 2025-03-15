package util

import (
	"regexp"
	"strings"
)

// ExtractCodeBlock extracts a code block from a string
func ExtractCodeBlock(input string) string {
	re := regexp.MustCompile("```[a-zA-Z0-9.]*\n([\\s\\S]+?)\n```")
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return input
}

// CodeBlockState represents the state of code block extraction
type CodeBlockState int

const (
	SearchingOpening CodeBlockState = iota
	Open
	Closed
)

// ExtractCodeBlockStream extracts code blocks from a stream
func ExtractCodeBlockStream(inputStream <-chan string) <-chan string {
	outputStream := make(chan string)

	go func() {
		defer close(outputStream)

		state := SearchingOpening
		buffer := strings.Builder{}
		openingRegex := regexp.MustCompile("```[^\n]*\n")
		potentialClosingRegex := regexp.MustCompile("\n`{0,2}$")

		for part := range inputStream {
			buffer.WriteString(part)
			bufStr := buffer.String()

			switch state {
			case SearchingOpening:
				match := openingRegex.FindStringIndex(bufStr)
				if match != nil {
					state = Open
					remainingContent := bufStr[match[1]:]
					buffer.Reset()
					buffer.WriteString(remainingContent)
				}

			case Open:
				if potentialClosingRegex.MatchString(bufStr) {
					continue
				}

				closePos := strings.Index(bufStr, "\n```")
				if closePos >= 0 {
					output := bufStr[:closePos]
					state = Closed
					buffer.Reset()
					outputStream <- output
					return
				}

				outputStream <- bufStr
				buffer.Reset()

			case Closed:
				return
			}
		}

		// If we've reached the end of the stream and still have content in the buffer
		if state != Closed && buffer.Len() > 0 {
			outputStream <- buffer.String()
		}
	}()

	return outputStream
}
