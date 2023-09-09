package utils

import (
	"context"

	"golang.design/x/clipboard"
)

var ctx = context.Background()

type CBReader <-chan string
type CBWriter chan string

const format = clipboard.FmtText

func InitClipboard() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	clipboard.Write(format, []byte(""))
}

func CleanupClipboard() {
	ctx.Done()
}

// GetClientClipboardIO retrieves channels for client clipboard I/O.
//
// This function sets up communication channels to perform clipboard I/O
// between the client and server. It returns two channels, one for writing
// data to the clipboard and another for receiving data from the clipboard.
//
// Returns:
//   - writer: A channel of type string for writing data to the clipboard.
//   - reader: A read-only channel of type string for reading data from
//     the clipboard.
func GetClientClipboardIO() (chan string, CBReader) {
	cb := getClipboardChan()

	writer := make(chan string)
	go writeToClipboard(AddMarkerToChannel(SERVER, writer))

	return writer, FilterChannelForMarker(CLIENT, cb)
}

// GetServerClipboardIO retrieves channels for server clipboard I/O.
//
// This function sets up communication channels to perform clipboard I/O
// between the client and server. It returns two channels, one for writing
// data to the clipboard and another for receiving data from the clipboard.
//
// Returns:
//   - writer: A channel of type string for writing data to the clipboard.
//   - reader: A read-only channel of type string for reading data from
//     the clipboard.
func GetServerClipboardIO() (CBWriter, CBReader) {
	cb := getClipboardChan()

	writer := make(chan string)
	go writeToClipboard(AddMarkerToChannel(CLIENT, writer))

	return writer, FilterChannelForMarker(SERVER, cb)
}

// GetControlClipboardIO retrieves channels for control/metadata clipboard I/O.
//
// This function sets up communication channels to perform clipboard I/O
// between the client and server. It returns two channels, one for writing
// data to the clipboard and another for receiving data from the clipboard.
//
// Returns:
//   - writer: A channel of type string for writing data to the clipboard.
//   - reader: A read-only channel of type string for reading data from
//     the clipboard.
func GetControlClipboardIO(marker Marker) (CBWriter, CBReader) {
	cb := getClipboardChan()

	writer := make(chan string)
	go writeToClipboard(AddMarkerToChannel(combineMarkers(CONTROL, getOppositeMarker(marker)), writer))

	return writer, FilterChannelForMarker(combineMarkers(CONTROL, marker), cb)
}

// returns a read-only string channel that includes the current contents
// of the clipboard and ongoing changes to the clipboard's contents
func getClipboardChan() <-chan string {
	changes := clipboard.Watch(ctx, format)
	output := make(chan string)

	go func() {
		defer close(output)
		for bytes := range changes {
			str := string(bytes)
			Debug("UPDATE: %s\n", str)
			output <- str
		}
	}()

	return output
}

// writes all contents of channel to the clipboard
func writeToClipboard(channel CBWriter) {
	for str := range channel {
		Debug("WRITE: %s\n", str)
		clipboard.Write(format, []byte(str))
	}
}
