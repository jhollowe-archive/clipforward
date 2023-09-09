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

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

// serverCmd represents the client command
var serverCmd = &cobra.Command{
	Use:   "server <dest> <port>",
	Short: "forwards traffic from the client to the 'port' port on 'dest'",
	// Long: ``,
	RunE: runServer,
	Args: cobra.ExactArgs(2),
}

func runServer(cmd *cobra.Command, args []string) error {
	fmt.Println("server called")

	utils.InitClipboard()

	proto, err := cmd.Flags().GetString("proto")
	cobra.CheckErr(err)

	dest := fmt.Sprintf("%s:%s", args[0], args[1])

	ctl_write, ctl_read := utils.GetControlClipboardIO(utils.SERVER)
	go handleServerCtl(ctl_write, ctl_read)

	for {
		handleServerConnection(proto, dest)
	}

	return nil
}

func handleServerConnection(proto string, dest string) {
	conn, err := net.Dial(proto, dest)
	cobra.CheckErr(err)
	defer conn.Close()

	s_write, s_read := utils.GetServerClipboardIO()

	go handleClientConnectionWrite(conn, s_read)
	handleClientConnectionRead(conn, s_write)
}

func handleServerConnectionRead(conn net.Conn, writer utils.CBWriter) {
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

func handleServerConnectionWrite(conn net.Conn, reader utils.CBReader) {
	for msg := range reader {
		utils.Debug("SERVER: " + msg)
		conn.Write([]byte(msg))
	}
}

func handleServerCtl(ctl_write utils.CBWriter, ctl_read utils.CBReader) {
	for msg := range ctl_read {
		utils.Debug("Server Ctl Handler: %s\n", msg)
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
