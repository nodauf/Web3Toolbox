package utilscmd

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

import (
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

// selectorCmd represents the selector command
var selectorCmd = &cobra.Command{
	Use:     "selector signature [...signature]",
	Short:   "A brief description of your command",
	Long:    ``,
	Example: `web3Toolbox utils selector balanceOf(address) transfer(address,uint256)`,
	Run: func(cmd *cobra.Command, args []string) {
		re := regexp.MustCompile(`\S+\(\S*\)`)
		for _, arg := range args {
			if !re.MatchString(arg) {
				fmt.Printf("%s is not a valid selector\n", arg)
				continue
			}
			fmt.Printf("%s: %s\n", arg, crypto.Keccak256Hash([]byte(arg)).String()[0:10])
		}
	},
}

func init() {
	UtilsCmd.AddCommand(selectorCmd)
	selectorCmd.Args = cobra.MinimumNArgs(1)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// selectorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// selectorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
