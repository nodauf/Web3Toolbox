package utilsDownloadSources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func init() {
	os.MkdirAll("./binSolc", 0755)
}

func DownloadSources(etherscanDomain, apiKey, contractAddress string) *GetSourceCodeEtherscanAPI {
	etherscanDownloadSourceAPI := fmt.Sprintf("%s/api?module=contract&action=getsourcecode&address=%s&apikey=%s", etherscanDomain, contractAddress, apiKey)

	res, err := http.Get(etherscanDownloadSourceAPI)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var sourceCodeEtherscanAPI GetSourceCodeEtherscanAPI
	err = json.Unmarshal(body, &sourceCodeEtherscanAPI)
	if err != nil {
		log.Fatal(err)
	}
	for i, result := range sourceCodeEtherscanAPI.Result {
		// If it's not json skip it
		if !strings.HasPrefix(result.SourceCodeString, "{{") {
			//fmt.Println("skip it")
			continue
		}
		decodedValue := result.SourceCodeString
		/*decodedValue, err := url.QueryUnescape(result.SourceCodeString)
		if err != nil {
			fmt.Println(result.SourceCodeString)
			log.Fatal("QueryUnescape" + err.Error())
		}*/
		// We have to remove the first { and last } as there is one more to be able to decode the json
		decodedValue = decodedValue[1 : len(decodedValue)-1]
		var sourceCode contractSourceCodeEtherscanAPI
		//fmt.Println(decodedValue)
		err = json.Unmarshal([]byte(decodedValue), &sourceCode)
		if err != nil {
			log.Fatal("unmarshal", err)
		}
		sourceCodeEtherscanAPI.Result[i].SourceCode = sourceCode
	}
	return &sourceCodeEtherscanAPI
}
