/*
Copyright © 2023 John Hollowell <jhollowe@johnhollowell.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"clipforward/utils"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client <port>",
	Short: "Opens a listening port <port> on the local machine which is forwarded to the server",
	// Long: ``,
	RunE: runClient,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(clientCmd)
}

func runClient(cmd *cobra.Command, args []string) error {
	fmt.Println("client called")

	utils.InitClipboard()

	proto, err := cmd.Flags().GetString("proto")
	cobra.CheckErr(err)

	// TODO allow setting the listen address to 0.0.0.0
	address := "127.0.0.1"

	ctl_write, ctl_read := utils.GetControlClipboardIO(utils.CLIENT)

	utils.Info("Attempting to find server...\n")
	for retries_left := 5; ; retries_left-- {
		if retries_left == 0 {
			errorAndExit(cmd, SERVER_NOT_FOUND, "Unable to find server. Please ensure the server is running and retry")
		}
		if !doPing(ctl_write, ctl_read) {
			utils.Info("Failed to find the server, retrying %d more times\n", retries_left)
			time.Sleep(CLIPBOARD_INTERVAL)
		} else {
			utils.Info("Server responded to clipboard ping\n")
			break
		}
	}

	// go handleClientCtl(ctl_write, ctl_read)

	// TODO use control message to check if there is a server running
	// TODO validate that server is using the same protocol
	// TODO validate the address is allowed

	// NETWORK LISTENER SECTION //
	listener, err := net.Listen(proto, address+":"+args[0])
	cobra.CheckErr(err)

	fmt.Printf("Listening for connections to %s\n", listener.Addr().String())

	for {
		fmt.Println("Waiting for new connection...")
		conn, err := listener.Accept()
		cobra.CheckErr(err)

		fmt.Printf("Got connection from %s\n", conn.RemoteAddr().String())

		// purposefully NOT a goroutine so that there is a single connection at a time
		// since the clipboard can only support one connection at a time
		handleClientConnection(conn)

		// TODO handle signals to nicely terminate
	}

	return nil
}

func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	c_write, c_read := utils.GetClientClipboardIO()

	go handleClientConnectionWrite(conn, c_read)
	handleClientConnectionRead(conn, c_write)
}

func handleClientConnectionRead(conn net.Conn, writer chan string) {
	buff := make([]byte, BUFF_SIZE)
	for {
		// read client request data
		count, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				utils.Error("failed to read data, err: %s", err)
			}
			utils.Info("Connection from %s closing\n", conn.RemoteAddr().String())
			return
		}
		utils.Debug("%d: '%s'\n", count, string(buff[:count]))

		writer <- string(buff[:count])

		// DEBUG just echo back
		// conn.Write(buff[:count])
	}
}

func handleClientConnectionWrite(conn net.Conn, reader <-chan string) {
	for str := range reader {
		if len(str) == 0 {
			return
		} else {
			utils.Debug("CLIENT: %s", str)
			conn.Write([]byte(str))
		}
	}
}

func handleClientCtl(ctl_write chan string, ctl_read <-chan string) {
	for msg := range ctl_read {
		msg_split := strings.SplitN(msg, utils.SEP, 1)
		cmd := msg_split[0]
		data := ""
		if len(msg_split) == 2 {
			data = msg_split[1]
		}

		switch cmd {
		case PING:
			ctl_write <- PONG
		case DISPLAY:
			utils.Info(data)
		case BYE:
			return
		}
	}
}
