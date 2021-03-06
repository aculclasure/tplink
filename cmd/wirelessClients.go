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
	"log"

	"github.com/spf13/cobra"
)

// wirelessClientsCmd represents the wirelessClients command
var wirelessClientsCmd = &cobra.Command{
	Use:   "wirelessClients",
	Short: "displays information about currently connected wireless clients",
	Long: `wirelessClients queries the wifi router to get the currently connected wireless clients and
prints out the IP address, MAC address, and host name (if known) for each wireless client.
`,
	Run: func(cmd *cobra.Command, args []string) {
		wireless, err := client.GetWirelessConnections()
		if err != nil {
			log.Fatalf("got error retrieving wireless connections (want a []*tplinkac9v1.Connection): %v", err)
		}
		if len(wireless) == 0 {
			fmt.Println("No wireless connections found")
		}
		fmt.Printf("%-15s%-22s%-15s\n", "IP_ADDRESS", "MAC_ADDRESS", "HOST_NAME")
		for _, w := range wireless {
			fmt.Printf("%-15s%-22s%-15s\n", w.IPAddress, w.MacAddress, w.Name)
		}
	},
}

func init() {
	listCmd.AddCommand(wirelessClientsCmd)
}
