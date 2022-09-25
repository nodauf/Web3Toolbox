/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utilscmd

import (
	"fmt"
	"github.com/nodauf/web3Toolbox/utilsCmd/reverseEvent"
	"github.com/spf13/cobra"
	"log"
)

//var online bool

// reverseEventCmd represents the reverseEvent command
var reverseEventCmd = &cobra.Command{
	Use:   "reverseEvent selector [...selector]",
	Short: "A brief description of your command",
	Long:  `The online flag will use less memory`,
	Run: func(cmd *cobra.Command, args []string) {
		if !online {

			for _, arg := range args {
				signatures, err := reverseEvent.ReverseEventFile(arg)
				if err != nil {
					log.Fatal(err)
				}
				for _, signature := range signatures {
					fmt.Printf("%s: %s\n", arg, signature)
				}
			}
		} else {
			for _, arg := range args {
				signatures, err := reverseEvent.ReverseEventOnline(arg)
				if err != nil {
					log.Fatal(err)
				}
				for _, signature := range signatures {
					fmt.Printf("%s: %s\n", arg, signature)
				}
			}
		}
	},
}

func init() {
	UtilsCmd.AddCommand(reverseEventCmd)
	reverseEventCmd.Args = cobra.MinimumNArgs(1)

	// Here you will define your flags and configuration settings.
	reverseEventCmd.Flags().BoolVarP(&online, "online", "o", false, "Use 4byteDirectory signatures")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reverseEventCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reverseEventCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
