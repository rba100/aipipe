// +build windows

package main

import (
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
)

const (
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
)

const newline = "\n"

// initConsole sets up the Windows console for proper UTF-8 output
func initConsole() {
	// Set stdout to UTF-8 mode
	handle := syscall.Handle(syscall.Stdout)
	
	var mode uint32
	procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	
	// Enable virtual terminal processing for ANSI escape sequences
	mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	
	procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
}
