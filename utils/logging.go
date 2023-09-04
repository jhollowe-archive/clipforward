package utils

import "fmt"

func info(format string, a ...any) {
	fmt.Printf(format, a...)
}

func debug(format string, a ...any) {
	fmt.Printf(format, a...)
}
