// +build !windows

package main

const newline = "\n"

// initConsole is a no-op on non-Windows platforms
func initConsole() {
	// Nothing to do on Unix-like systems
}
