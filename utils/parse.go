package utils

type Marker string

const (
	NONE             Marker = "X"
	CLIENT           Marker = "C"
	SERVER           Marker = "S"
	CONTROL          Marker = "M" // for Metadata
	marker_separator string = ";"
)

func FilterChannelForMarker(marker Marker, channel <-chan string) <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		for str := range channel {
			if len(str) > 0 {
				Debug(str)
				// payload_str, has_prefix := strings.CutPrefix(str, string(marker)+marker_separator)
				// if has_prefix {
				// 	output <- payload_str
				// }
				output <- str
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
			output <- string(marker) + marker_separator + str
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
