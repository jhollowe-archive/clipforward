package utils

import "strings"

func FilterChannelForMarker(marker Marker, channel <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		for raw_str := range channel {
			if len(raw_str) > 0 {
				payload_str, has_prefix := strings.CutPrefix(raw_str, string(marker)+SEP)
				if has_prefix {
					output <- payload_str
				}
			}
		}
	}()

	return output
}

func AddMarkerToChannel(marker Marker, channel <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		for str := range channel {
			output <- string(marker) + SEP + str
		}
	}()

	return output
}

func getOppositeMarker(marker Marker) Marker {
	switch marker {
	case CLIENT:
		return SERVER
	case SERVER:
		return CLIENT
	default:
		return NONE
	}
}

func combineMarkers(marks ...Marker) Marker {
	var combined Marker
	for _, m := range marks {
		combined += m
	}
	return combined
}
