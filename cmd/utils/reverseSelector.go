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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/nodauf/web3Toolbox/embeddedFiles"
	"github.com/spf13/cobra"
)

type signature struct {
	ID            int    `json:"id"`
	TextSignature string `json:"text_signature"`
	HexSignature  string `json:"hex_signature"`
}

type response4Bytes struct {
	Results []signature `json:"results"`
}

var online bool

// reverseSelectorCmd represents the reverseSelector command
var reverseSelectorCmd = &cobra.Command{
	Use:   "reverseSelector selector [...selector]",
	Short: "A brief description of your command",
	Long:  `The online flag will use less memory`,
	Run: func(cmd *cobra.Command, args []string) {
		if !online {
			var allSignatures []signature
			err := json.Unmarshal(embeddedFiles.Get4bytesSignatures(), &allSignatures)
			if err != nil {
				log.Fatalf("Fail to unmarshal: %s", err.Error())
				return
			}
			// Create index for signatures to avoid complexity O(n²) -> O(2n)
			var indexSignatures = make(map[string][]string, len(allSignatures))
			for _, signature := range allSignatures {
				indexSignatures[signature.HexSignature] = append(indexSignatures[signature.HexSignature], signature.TextSignature)

			}
			for _, arg := range args {
				if textSignature, ok := indexSignatures[arg]; ok {
					fmt.Printf("%s: %s\n", arg, strings.Join(textSignature, "\n"+arg+": "))
					// We continue, one selector can match multiple signatures
					//break
				}
			}
		} else {
			for _, arg := range args {
				var baseURL = "https://www.4byte.directory/api/v1/signatures/?hex_signature=%s"
				url := fmt.Sprintf(baseURL, arg)
				resp, err := http.Get(url)
				if err != nil {
					fmt.Printf("Fail to get page at %s: %s\n", url, err.Error())
					continue
				}
				if resp.StatusCode != 200 {
					resp.Body.Close()
					fmt.Printf("Fail to get page at %s status code %d:\n", url, resp.StatusCode)
					continue
				}
				var resp4Bytes response4Bytes
				err = json.NewDecoder(resp.Body).Decode(&resp4Bytes)
				resp.Body.Close()
				if err != nil {
					fmt.Printf("Fail to decode response at %s: %s\n", url, err.Error())
					continue
				}
				for _, signature := range resp4Bytes.Results {
					fmt.Printf("%s: %s\n", arg, signature.TextSignature)
				}
			}
		}

	},
}

func init() {
	UtilsCmd.AddCommand(reverseSelectorCmd)
	reverseSelectorCmd.Args = cobra.MinimumNArgs(1)

	// Here you will define your flags and configuration settings.
	reverseSelectorCmd.Flags().BoolVarP(&online, "online", "o", false, "Use 4byteDirectory signatures")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reverseSelectorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reverseSelectorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
