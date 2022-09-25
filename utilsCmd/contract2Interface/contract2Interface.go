package contract2Interface

import (
	"fmt"
	"github.com/nodauf/solc-go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func GetSolidityVersion(contractContent []byte) string {
	firstLine := strings.Split(string(contractContent), "\n")[0]
	re := regexp.MustCompile(`^pragma solidity .?(\d{1,2}\.\d{1,2}\.\d{1,2});.*`)
	match := re.FindStringSubmatch(firstLine)
	if len(match) == 2 {
		return match[1]
	}
	secondLine := strings.Split(string(contractContent), "\n")[1]
	return re.FindStringSubmatch(secondLine)[1]

}

func DownloadSolcJS(solcVersion string) string {

	err := os.MkdirAll("./solc-bin", 0755)
	if err != nil {
		log.Fatal("Fail to create directory solc-bin", err)
	}

	res, err := http.Get("https://github.com/ethereum/solc-bin/raw/gh-pages/bin/list.txt")
	if err != nil {
		log.Fatal("http.get solc", err)
	}

	listSolcVersion, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal("readAll solc", err)
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

			// Create the file
			out, err := os.Create("./solc-bin/" + line)
			if err != nil {
				resp.Body.Close()
				log.Fatal(err)
			}
			os.Chmod("./solc-bin/"+line, 0755)
			// Write the body to file
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				resp.Body.Close()
				log.Fatalf("fail to write to file %s: %s", "./solc-bin/"+line, err.Error())

			}
			resp.Body.Close()
			out.Close()
			solcBin = "./solc-bin/" + line
			break
		}
	}
	return solcBin
}

func InterfaceFromAST(ast solc.ASTsolc) string {
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
					var visibility string
					if contractNode.Visibility == "public" {
						visibility = "external"
					} else {
						visibility = contractNode.Visibility
					}
					interfaceString += fmt.Sprintf(") %s returns(%s);\n", visibility, returnType)

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
					var visibility string
					if contractNode.Visibility == "public" {
						visibility = "external"
					} else {
						visibility = contractNode.Visibility
					}
					interfaceString += fmt.Sprintf(") %s", visibility)
					if len(contractNode.ReturnParameters.Parameters) > 0 {
						interfaceString += fmt.Sprint(" returns (")
						for l, returnParameter := range contractNode.ReturnParameters.Parameters {
							interfaceString += fmt.Sprintf("%s", returnParameter.TypeDescriptions.TypeString)
							if l != len(contractNode.ReturnParameters.Parameters)-1 {
								interfaceString += fmt.Sprint(", ")
							}
						}
						interfaceString += fmt.Sprint(")")
					}
					interfaceString += fmt.Sprint(";\n")
				}
				if j == len(node.Nodes)-1 {

					interfaceString += "}"
				}
			}
		}

	}
	return interfaceString
}
