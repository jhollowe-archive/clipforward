package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

func doPing(ctl_write chan string, ctl_read <-chan string) bool {

	// Attempt a ping in a goroutine
	ping_result := make(chan bool, 1)
	go func() {
		ctl_write <- CLEAR
		time.Sleep(CLIPBOARD_INTERVAL)
		ctl_write <- PING

		// wait for a response and return false if timed out
		select {
		case msg := <-ctl_read:
			ping_result <- (msg == PONG)
		case <-time.After(CLIPBOARD_INTERVAL * 3):
			ping_result <- false
		}
	}()

	return <-ping_result
}

func errorAndExit(cmd *cobra.Command, exit_code int, str string) {
	cmd.PrintErrln(str)
	os.Exit(exit_code)
}
