/*
Copyright © 2020 Andrew Culclasure

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
	"github.com/aculclasure/tplink/archerc9v1"
	"github.com/spf13/cobra"
	"log"
)

// rebootCmd represents the reboot command
var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "reboots the router",
	Long:  `reboot restarts the router!`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		client, err = archerc9v1.New(userName, password, url, nil, nil)
		if err != nil {
			log.Fatalf("got error trying to create new tplinkac9v1.Client: %s", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := client.Reboot(); err != nil {
			log.Fatalf("got error rebooting the router (want a response that the router rebooted successfuly): %v", err)
		}
		fmt.Println("router rebooted!")
	},
}

func init() {
	rootCmd.AddCommand(rebootCmd)
	rebootCmd.Flags().StringVar(&url, "url", "http://192.168.168.1", "router URL (required)")
	rebootCmd.MarkFlagRequired("url")
	rebootCmd.Flags().StringVarP(&userName, "user", "U", "admin", "router admin user name (required)")
	rebootCmd.MarkFlagRequired("user")
	rebootCmd.Flags().StringVarP(&password, "password", "P", "admin", "router admin password (required)")
	rebootCmd.MarkFlagRequired("password")
}
