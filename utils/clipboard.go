package utils

import (
	"context"

	"golang.design/x/clipboard"
)

var ctx = context.Background()

const format = clipboard.FmtText

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
func GetClientClipboardIO() (chan string, <-chan string) {
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
func GetServerClipboardIO() (chan string, <-chan string) {
	cb := getClipboardChan()

	writer := make(chan string)
	go writeToClipboard(AddMarkerToChannel(CLIENT, writer))

	return writer, FilterChannelForMarker(SERVER, cb)
}

// returns a read-only string channel that includes the current contents
// of the clipboard and ongoing changes to the clipboard's contents
func getClipboardChan() <-chan string {
	changes := clipboard.Watch(ctx, format)
	output := make(chan string)

	output <- string(clipboard.Read(format))

	go func() {
		defer close(output)
		for bytes := range changes {
			str := string(bytes)
			output <- str
		}
	}()

	return output
}

// writes all contents of channel to the clipboard
func writeToClipboard(channel <-chan string) {
	for str := range channel {
		clipboard.Write(format, []byte(str))
	}
}
