package utils

import (
	"fmt"
	"os"
)

func Info(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Debug(format string, a ...any) {
	fmt.Printf(format, a...)
}

func Error(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}
