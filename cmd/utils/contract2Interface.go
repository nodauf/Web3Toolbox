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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nodauf/solc-go"
	"github.com/spf13/cobra"
)

var contratFile string

// contract2InterfaceCmd represents the contract2Interface command
var contract2InterfaceCmd = &cobra.Command{
	Use:   "contract2Interface",
	Short: "Convert a solidity contract to an inteface",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		contractData, err := os.ReadFile(contratFile)
		if err != nil {
			log.Fatal("Fail to open file", err)
		}
		downloadedFile := downloadSolcJS(getSolidityVersion(contractData))
		compiler, err := solc.NewFromFile(downloadedFile)
		if err != nil {
			log.Fatal(err)
		}
		input := &solc.Input{
			Language: "Solidity",
			Sources: map[string]solc.SourceIn{
				"target.sol": solc.SourceIn{Content: string(contractData)},
			},
			Settings: solc.Settings{
				EVMVersion: "byzantium",
				OutputSelection: map[string]map[string][]string{
					"*": map[string][]string{
						"*": []string{
							//"abi",
							//"evm.bytecode.object",
							//"evm.bytecode.sourceMap",
							//"evm.deployedBytecode.object",
							//"evm.deployedBytecode.sourceMap",
							//"evm.methodIdentifiers",
						},
						"": []string{
							"ast",
						},
					},
				},
			},
		}
		output, err := compiler.Compile(input)
		if err != nil {
			log.Fatal(err)
		}
		//x, _ := json.MarshalIndent(output.Sources["target.sol"].AST, "", "	")
		//fmt.Printf("%s ", x)
		fmt.Println(interfaceFromAST(output.Sources["target.sol"].AST))
		//var interface string

	},
}

func init() {
	UtilsCmd.AddCommand(contract2InterfaceCmd)
	contract2InterfaceCmd.Flags().StringVarP(&contratFile, "file", "f", "", "contract file")
	contract2InterfaceCmd.MarkFlagRequired("file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// contract2InterfaceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// contract2InterfaceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getSolidityVersion(contractContent []byte) string {
	firstLine := strings.Split(string(contractContent), "\n")[0]
	re, _ := regexp.Compile(`^pragma solidity .?([0-9]{1,2}\.[0-9]{1,2}\.[0-9]{1,2});.*`)
	match := re.FindStringSubmatch(firstLine)
	if len(match) == 2 {
		return match[1]
	}
	secondLine := strings.Split(string(contractContent), "\n")[1]
	return re.FindStringSubmatch(secondLine)[1]

}

func downloadSolcJS(solcVersion string) string {

	os.MkdirAll("./solc-bin", 0755)

	res, err := http.Get("https://github.com/ethereum/solc-bin/raw/gh-pages/bin/list.txt")
	if err != nil {
		log.Fatal("http.get solc", err)
	}

	listSolcVersion, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("readall solc", err)
	}
	var solcBin string
	for _, line := range strings.Split(string(listSolcVersion), "\n") {
		if strings.Contains(line, solcVersion+"+") {
			// Download the file
			// Get the data
			resp, err := http.Get("https://github.com/ethereum/solc-bin/raw/gh-pages/bin/" + line)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			// Create the file
			out, err := os.Create("./solc-bin/" + line)
			if err != nil {
				log.Fatal(err)
			}
			defer out.Close()
			os.Chmod("./solc-bin/"+line, 0755)
			// Write the body to file
			_, err = io.Copy(out, resp.Body)
			solcBin = "./solc-bin/" + line
			break
		}
	}
	return solcBin
}

func interfaceFromAST(ast solc.ASTsolc) string {
	var interfaceString string
	for _, node := range ast.Nodes {
		if node.NodeType == "PragmaDirective" {
			interfaceString += fmt.Sprintf("pragma %s %s;\n", node.Literals[0], strings.Join(node.Literals[1:], ""))
		}
		if node.NodeType == "ContractDefinition" {
			interfaceString += fmt.Sprintf("interface %s {\n", node.Name)
			for j, contractNode := range node.Nodes {

				// Variables
				if contractNode.NodeType == "VariableDeclaration" && contractNode.Name != "" {
					var returnType string
					interfaceString += fmt.Sprintf("\tfunction %s(", contractNode.Name)
					if contractNode.TypeName.NodeType == "Mapping" {
						interfaceString += fmt.Sprintf("%s key", contractNode.TypeName.KeyType.Name)
						returnType = contractNode.TypeName.ValueType.Name
					} else {
						returnType = contractNode.TypeDescriptions.TypeString

					}
					interfaceString += fmt.Sprintf(") %s returns(%s);\n", contractNode.Visibility, returnType)

				}

				// Functions
				if contractNode.NodeType == "FunctionDefinition" && contractNode.Name != "" {
					interfaceString += fmt.Sprintf("\tfunction %s(", contractNode.Name)
					// Parse parameters
					for k, parameter := range contractNode.Parameters.Parameters {
						interfaceString += fmt.Sprintf("%s ", parameter.TypeDescriptions.TypeString)
						if parameter.StorageLocation != "default" {
							interfaceString += fmt.Sprintf("%s ", parameter.StorageLocation)
						}
						interfaceString += fmt.Sprintf("%s", parameter.Name)

						if k != len(contractNode.Parameters.Parameters)-1 {
							interfaceString += fmt.Sprintf(", ")
						}
					}
					interfaceString += fmt.Sprintf(") %s", contractNode.Visibility)
					if len(contractNode.ReturnParameters.Parameters) > 0 {
						interfaceString += fmt.Sprintf(" returns (")
						for l, returnParameter := range contractNode.ReturnParameters.Parameters {
							interfaceString += fmt.Sprintf("%s", returnParameter.TypeDescriptions.TypeString)
							if l != len(contractNode.ReturnParameters.Parameters)-1 {
								interfaceString += fmt.Sprintf(", ")
							}
						}
						interfaceString += fmt.Sprintf(")")
					}
					interfaceString += fmt.Sprintf(";\n")
				}
				if j == len(node.Nodes)-1 {

					interfaceString += "}"
				}
			}
		}

	}
	return interfaceString
}
