package display

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// ColorMode represents the color mode supported by the terminal
type ColorMode int

const (
	// SimpleColorMode represents basic 16-color ANSI support
	SimpleColorMode ColorMode = iota
	// Color256Mode represents 256-color ANSI support
	Color256Mode
	// TrueColorMode represents 24-bit true color support
	TrueColorMode
)

// Format constants for text formatting
const (
	// ResetFormat resets all formatting
	ResetFormat = "\033[0m"
	// BoldFormat makes text bold
	BoldFormat = "\033[1m"
	// DimFormat makes text dim
	DimFormat = "\033[2m"
	// ItalicFormat makes text italic (not supported in all terminals)
	ItalicFormat = "\033[3m"
	// UnderlineFormat makes text underlined
	UnderlineFormat = "\033[4m"
)

// Simple color constants (16 colors)
const (
	// Foreground colors
	BlackFg   = "\033[30m"
	RedFg     = "\033[31m"
	GreenFg   = "\033[32m"
	YellowFg  = "\033[33m"
	BlueFg    = "\033[34m"
	MagentaFg = "\033[35m"
	CyanFg    = "\033[36m"
	WhiteFg   = "\033[37m"

	// Bright foreground colors
	BrightBlackFg   = "\033[90m"
	BrightRedFg     = "\033[91m"
	BrightGreenFg   = "\033[92m"
	BrightYellowFg  = "\033[93m"
	BrightBlueFg    = "\033[94m"
	BrightMagentaFg = "\033[95m"
	BrightCyanFg    = "\033[96m"
	BrightWhiteFg   = "\033[97m"

	// Background colors
	BlackBg   = "\033[40m"
	RedBg     = "\033[41m"
	GreenBg   = "\033[42m"
	YellowBg  = "\033[43m"
	BlueBg    = "\033[44m"
	MagentaBg = "\033[45m"
	CyanBg    = "\033[46m"
	WhiteBg   = "\033[47m"
)

// Color mappings for syntax highlighting and markdown formatting
var (
	// Token type colors for syntax highlighting (will be initialized in InitializeColors)
	TokenKeywordColor    string
	TokenIdentifierColor string
	TokenLiteralColor    string
	TokenCommentColor    string
	TokenOtherColor      string

	// Markdown formatting colors (will be initialized in InitializeColors)
	MdHeaderColor     string
	MdCodeBlockColor  string
	MdInlineCodeColor string
	MdBlockQuoteColor string
	MdListMarkerColor string
	MdEmphasisColor   string
	MdHorizontalColor string
	MdNormalTextColor string
)

// GetColorMode detects the terminal's color capabilities
func GetColorMode() ColorMode {
	// Check for true color support
	if os.Getenv("COLORTERM") == "truecolor" || os.Getenv("COLORTERM") == "24bit" {
		return TrueColorMode
	}

	// Check for 256 color support
	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return Color256Mode
	}

	// Default to simple color mode
	return SimpleColorMode
}

// IsBoldSupported checks if bold formatting is supported
func IsBoldSupported() bool {
	// Windows Terminal supports bold, but cmd.exe doesn't
	return runtime.GOOS != "windows" || IsWindowsTerminal()
}

// Get256Color returns a 256-color ANSI code for the given color index
func Get256Color(index int, isForeground bool) string {
	if isForeground {
		return fmt.Sprintf("\033[38;5;%dm", index)
	}
	return fmt.Sprintf("\033[48;5;%dm", index)
}

// GetRGBColor returns a true color ANSI code for the given RGB values
func GetRGBColor(r, g, b int, isForeground bool) string {
	if isForeground {
		return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
	}
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

// IsWindowsTerminal checks if the current terminal is Windows Terminal
func IsWindowsTerminal() bool {
	return os.Getenv("WT_SESSION") != ""
}

// InitializeColors sets up colors based on terminal capabilities
func InitializeColors() {
	mode := GetColorMode()

	// Set default colors (16-color mode)
	TokenKeywordColor = MagentaFg
	TokenIdentifierColor = WhiteFg
	TokenLiteralColor = GreenFg
	TokenCommentColor = BrightBlackFg
	TokenOtherColor = CyanFg

	MdHeaderColor = BoldFormat + YellowFg
	MdCodeBlockColor = CyanFg
	MdInlineCodeColor = CyanFg
	MdBlockQuoteColor = BlueFg
	MdListMarkerColor = BlueFg
	MdEmphasisColor = YellowFg + DimFormat
	MdHorizontalColor = YellowFg
	MdNormalTextColor = WhiteFg

	// Windows Terminal has good color support even if TERM doesn't indicate it
	if IsWindowsTerminal() {
		mode = Color256Mode
	}

	// If we have 256 color support, use more vibrant colors
	if mode == Color256Mode || mode == TrueColorMode {
		TokenKeywordColor = "\033[38;5;171m"    // Bright purple
		TokenIdentifierColor = "\033[38;5;252m" // Light gray
		TokenLiteralColor = "\033[38;5;114m"    // Light green
		TokenCommentColor = "\033[38;5;245m"    // Medium gray
		TokenOtherColor = "\033[38;5;81m"       // Light cyan

		MdHeaderColor = BoldFormat + "\033[38;5;220m" // Gold
		MdCodeBlockColor = "\033[38;5;81m"            // Light cyan
		MdInlineCodeColor = "\033[38;5;81m"           // Light cyan
		MdBlockQuoteColor = "\033[38;5;75m"           // Medium blue
		MdListMarkerColor = "\033[38;5;75m"           // Medium blue
		MdEmphasisColor = "\033[38;5;222m"            // Light gold
		MdHorizontalColor = "\033[38;5;220m"          // Gold
		MdNormalTextColor = "\033[38;5;252m"          // Light gray
	}
}
