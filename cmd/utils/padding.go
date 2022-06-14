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

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var hex bool
var size int
var direction string

// paddingCmd represents the padding command
var paddingCmd = &cobra.Command{
	Use:   "padding [Flags] string",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if direction != "l" && direction != "r" {
			return fmt.Errorf("direction must be l or r")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var value []byte
		if hex {
			value = common.FromHex(args[0])
		} else {
			value = []byte(args[0])
		}
		if direction == "l" {
			fmt.Printf("padded string in hex: %s\n", common.Bytes2Hex(common.LeftPadBytes(value, 32)))
			fmt.Printf("keccack256 hash of previous hex: %s\n", crypto.Keccak256Hash(common.LeftPadBytes(value, 32)))
		} else {
			fmt.Printf("padded string in hex: %s\n", common.Bytes2Hex(common.RightPadBytes(value, size)))
			fmt.Printf("keccack256 hash of previous hex: %s\n", crypto.Keccak256Hash(common.RightPadBytes(value, size)))
		}
	},
}

func init() {
	UtilsCmd.AddCommand(paddingCmd)

	paddingCmd.Args = cobra.ExactArgs(1)
	paddingCmd.Flags().StringVarP(&direction, "direction", "d", "l", "direction of padding (r for right or l for left)")
	paddingCmd.Flags().IntVarP(&size, "size", "s", 32, "size of the padding")
	paddingCmd.Flags().BoolVar(&hex, "hex", false, "Interpret the string as hex")
}
