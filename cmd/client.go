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
	"fmt"
	"net"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client <port>",
	Short: "Opens a listening port <port> on the local machine which is forwarded to the server",
	// Long: ``,
	RunE: runClient,
	Args: cobra.ExactArgs(1),
}

func runClient(cmd *cobra.Command, args []string) error {
	fmt.Println("client called")

	// TODO use control message to check if there is a server running

	// TODO validate that server is using the same protocol

	// TODO validate the address is allowed

	listener, err := net.Listen(viper.GetString("proto"), viper.GetString("address"))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			// TODO handle error
		}
		go handleConnection(conn)
	}

	return nil
}

func handleConnection(conn net.Conn) {
	var input []byte
	count, err := conn.Read(input)
	if err != nil {
		// TODO handle error
	}
	fmt.Printf("%d: %s", count, string(input))

	// DEBUG just echo back
	conn.Write(input)
}

func init() {
	rootCmd.AddCommand(clientCmd)

}
