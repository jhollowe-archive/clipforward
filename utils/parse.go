package utils

import (
	"strings"
)

const (
	CLIENT           = "C"
	SERVER           = "S"
	CONTROL          = "M" // for Metadata
	marker_separator = ";"
)

func FilterChannelForMarker(marker string, channel <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		for str := range channel {
			if len(str) > 0 && strings.HasPrefix(str, marker+marker_separator) {
				output <- str
			}
		}
	}()

	return output
}

func AddMarkerToChannel(marker string, channel <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		for str := range channel {
			output <- marker + marker_separator + str
		}
	}()

	return output
}
