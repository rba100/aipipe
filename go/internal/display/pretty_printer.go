package display

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
)

// PrintState represents the current state of the pretty printer
type PrintState int

const (
	Normal PrintState = iota
	InCodeBlock
)

// PrettyPrinter handles pretty printing of markdown text
type PrettyPrinter struct {
	originalColor       int
	isBoldSupported     bool
	lineBuffer          strings.Builder
	currentState        PrintState
	headerRegex         *regexp.Regexp
	inlineCodeRegex     *regexp.Regexp
	codeBlockStartRegex *regexp.Regexp
	codeBlockEndRegex   *regexp.Regexp
	numberedListRegex   *regexp.Regexp
	unorderedListRegex  *regexp.Regexp
	emphasisRegex       *regexp.Regexp
	blockQuoteRegex     *regexp.Regexp
	horizontalRuleRegex *regexp.Regexp
}

// ANSI color codes
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	Grey       = "\033[90m"
	DarkYellow = "\033[33;2m"
)

// NewPrettyPrinter creates a new pretty printer
func NewPrettyPrinter() *PrettyPrinter {
	isBoldSupported := runtime.GOOS != "windows" || os.Getenv("WT_SESSION") != ""

	p := &PrettyPrinter{
		originalColor:   0, // Not used in Go implementation
		isBoldSupported: isBoldSupported,
		currentState:    Normal,
		lineBuffer:      strings.Builder{},
	}

	p.headerRegex = regexp.MustCompile(`^#{1,6}\s+.*$`)
	p.inlineCodeRegex = regexp.MustCompile("\x60[^\x60\n]+\x60")
	p.codeBlockStartRegex = regexp.MustCompile(`^\s*\x60\x60\x60`)
	p.codeBlockEndRegex = regexp.MustCompile(`^\s*\x60\x60\x60\s*$`)
	p.numberedListRegex = regexp.MustCompile(`^(\s*)(\d+\.)\s+(.*)$`)
	p.unorderedListRegex = regexp.MustCompile(`^(\s*)([-*])\s+(.*)$`)
	p.emphasisRegex = regexp.MustCompile(`(\*\*\*|\*\*|__)([^*_]+)(\*\*\*|\*\*|__)|(\*|_)([^*_]+)(\*|_)`)
	p.blockQuoteRegex = regexp.MustCompile(`^(\s*)((?:>\s*)+)(.*)$`)
	p.horizontalRuleRegex = regexp.MustCompile(`^(\s*)([-*_])([-*_])([-*_])+\s*$`)

	return p
}

// Close cleans up the pretty printer
func (p *PrettyPrinter) Close() {
	fmt.Print(Reset)
}

// Flush prints any remaining content in the line buffer
func (p *PrettyPrinter) Flush() {
	if p.lineBuffer.Len() > 0 {
		var line string = p.lineBuffer.String()
		p.processLine(line)
		p.lineBuffer.Reset()
		if !strings.HasSuffix(line, "\n") {
			fmt.Println()
		}
	}
}

// Print prints the text with pretty formatting
func (p *PrettyPrinter) Print(text string) {
	if len(text) == 0 {
		return
	}

	if !strings.Contains(text, "\n") {
		p.lineBuffer.WriteString(text)
		return
	}

	isTerminated := strings.HasSuffix(text, "\n")
	if isTerminated {
		text = text[:len(text)-1]
	}

	lines := strings.Split(text, "\n")
	for i, line := range lines {
		isLastLine := i == len(lines)-1

		if p.lineBuffer.Len() > 0 {
			line = p.lineBuffer.String() + line
			p.lineBuffer.Reset()
		}

		if isLastLine && !isTerminated {
			p.lineBuffer.WriteString(line)
			return
		}

		p.processLine(line)
		if !isLastLine {
			fmt.Println()
		}
	}

	if isTerminated {
		fmt.Println()
	}
}

// processLine processes a single line of text
func (p *PrettyPrinter) processLine(line string) {
	if strings.Contains(line, "\r") {
		line = strings.ReplaceAll(line, "\r", "")
	}

	if p.currentState == Normal {
		if p.codeBlockStartRegex.MatchString(line) {
			fmt.Print(Cyan)
			fmt.Print(line)
			p.currentState = InCodeBlock
			return
		}

		p.processNormalLine(line)
	} else { // InCodeBlock
		if p.codeBlockEndRegex.MatchString(line) {
			fmt.Print(Cyan)
			fmt.Print(line)
			p.currentState = Normal
			return
		}

		fmt.Print(Cyan)
		fmt.Print(line)
	}
}

// processNormalLine processes a line in normal (non-code-block) state
func (p *PrettyPrinter) processNormalLine(line string) {
	if p.headerRegex.MatchString(line) {
		p.printHeader(line)
		return
	}

	if p.horizontalRuleRegex.MatchString(line) {
		p.printHorizontalRule(line)
		return
	}

	if p.blockQuoteRegex.MatchString(line) {
		p.printBlockQuote(line)
		return
	}

	if p.numberedListRegex.MatchString(line) {
		p.printNumberedList(line)
		return
	}

	if p.unorderedListRegex.MatchString(line) {
		p.printUnorderedList(line)
		return
	}

	p.printFormattedText(line)
}

// printHeader prints a header line
func (p *PrettyPrinter) printHeader(line string) {
	fmt.Print(Bold + Yellow)
	fmt.Print(line)
	fmt.Print(Reset)
}

// printHorizontalRule prints a horizontal rule
func (p *PrettyPrinter) printHorizontalRule(line string) {
	fmt.Print(Yellow)
	fmt.Print(line)
	fmt.Print(Reset)
}

// printBlockQuote prints a block quote
func (p *PrettyPrinter) printBlockQuote(line string) {
	matches := p.blockQuoteRegex.FindStringSubmatch(line)
	if len(matches) >= 4 {
		indentation := matches[1]
		quote := matches[2]
		content := matches[3]

		fmt.Print(indentation)
		fmt.Print(Blue)
		fmt.Print(quote)
		fmt.Print(Reset)
		p.printFormattedText(content)
	}
}

// printNumberedList prints a numbered list item
func (p *PrettyPrinter) printNumberedList(line string) {
	matches := p.numberedListRegex.FindStringSubmatch(line)
	if len(matches) >= 4 {
		indentation := matches[1]
		number := matches[2]
		content := matches[3]

		fmt.Print(indentation)
		fmt.Print(Blue)
		fmt.Print(number)
		fmt.Print(Reset)
		fmt.Print(" ")
		p.printFormattedText(content)
	}
}

// printUnorderedList prints an unordered list item
func (p *PrettyPrinter) printUnorderedList(line string) {
	matches := p.unorderedListRegex.FindStringSubmatch(line)
	if len(matches) >= 4 {
		indentation := matches[1]
		bullet := matches[2]
		content := matches[3]

		fmt.Print(indentation)
		fmt.Print(Blue)
		fmt.Print(bullet)
		fmt.Print(Reset)
		fmt.Print(" ")
		p.printFormattedText(content)
	}
}

// printFormattedText prints text with inline formatting
func (p *PrettyPrinter) printFormattedText(line string) {
	lastIndex := 0
	inlineCodeMatches := p.inlineCodeRegex.FindAllStringIndex(line, -1)
	emphasisMatches := p.emphasisRegex.FindAllStringIndex(line, -1)

	// Combine and sort all matches by index
	type match struct {
		index  int
		length int
		typ    string
	}

	allMatches := []match{}

	for _, m := range inlineCodeMatches {
		allMatches = append(allMatches, match{
			index:  m[0],
			length: m[1] - m[0],
			typ:    "code",
		})
	}

	for _, m := range emphasisMatches {
		allMatches = append(allMatches, match{
			index:  m[0],
			length: m[1] - m[0],
			typ:    "emphasis",
		})
	}

	// Sort matches by index
	for i := 0; i < len(allMatches); i++ {
		for j := i + 1; j < len(allMatches); j++ {
			if allMatches[i].index > allMatches[j].index {
				allMatches[i], allMatches[j] = allMatches[j], allMatches[i]
			}
		}
	}

	for _, m := range allMatches {
		// Print text before the match
		if m.index > lastIndex {
			fmt.Print(White)
			fmt.Print(line[lastIndex:m.index])
		}

		// Print the match with appropriate formatting
		if m.typ == "code" {
			fmt.Print(Cyan)
			fmt.Print(line[m.index : m.index+m.length])
		} else if m.typ == "emphasis" {
			fmt.Print(DarkYellow)
			if p.isBoldSupported {
				fmt.Print(Bold)
				fmt.Print(line[m.index : m.index+m.length])
				fmt.Print(Reset + White) // Reset bold but keep color
			} else {
				fmt.Print(line[m.index : m.index+m.length])
			}
		}

		lastIndex = m.index + m.length
	}

	// Print remaining text
	if lastIndex < len(line) {
		fmt.Print(White)
		fmt.Print(line[lastIndex:])
	}

	fmt.Print(Reset)
}
