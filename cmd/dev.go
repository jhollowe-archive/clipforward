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
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: runDev,
}

func runDev(cmd *cobra.Command, args []string) error {
	fmt.Println("dev called")

	utils.InitClipboard()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	fmt.Println(string(clipboard.Read(clipboard.FmtText)))

	c_write, c_read := utils.GetClientClipboardIO()
	s_write, s_read := utils.GetServerClipboardIO()
	go func(ctx context.Context) {
		c_write <- "hello world from the client"
		time.Sleep(1100 * time.Millisecond)
		s_write <- "hello world from the server"
	}(ctx)
	go readUntilTimeout(c_read, "CLIENT", 120*time.Second)
	readUntilTimeout(s_read, "SERVER", 120*time.Second)

	return nil
}

func init() {
	rootCmd.AddCommand(devCmd)
}

func readUntilTimeout(input <-chan string, prefix string, timeout time.Duration) {
	timer := time.NewTimer(timeout)

	for {
		select {
		case value := <-input:
			fmt.Printf("%s: %s\n", prefix, value)
		case <-timer.C:
			return
		}
	}
}
