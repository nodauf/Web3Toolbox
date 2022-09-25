package ethcmd

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
	"github.com/spf13/cobra"
	solidityParser "github.com/umbracle/solidity-parser-go"
	"log"
	"os"
)

var contractPath string

// solidityParserCmd represents the solidityParser command
var solidityParserCmd = &cobra.Command{
	Use:   "solidityParser",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		body, err := os.ReadFile(contractPath)
		if err != nil {
			log.Fatalln(err)
		}
		parser := solidityParser.Parse(string(body))
		json, err := parser.Json()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(json)
	},
}

func init() {
	EthCmd.AddCommand(solidityParserCmd)
	solidityParserCmd.Flags().StringVarP(&contractPath, "contractPath", "c", "", "Contract path")
	//solidityParserCmd.Args = cobra.MinimumNArgs(1)
	solidityParserCmd.MarkFlagRequired("contractPath")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// solidityParserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// solidityParserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
