package reverseEvent

import (
	"encoding/json"
	"fmt"
	"github.com/nodauf/web3Toolbox/embeddedFiles"
	"net/http"
)

func ReverseEventOnline(Event string) ([]string, error) {
	var baseURL = "https://www.4byte.directory/api/v1/event-signatures/?hex_signature=%s"
	url := fmt.Sprintf(baseURL, Event)
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Printf("Fail to get page at %s: %s\n", url, err.Error())
		return nil, fmt.Errorf("Fail to get page at %s: %s", url, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//fmt.Printf("Fail to get page at %s status code %d:\n", url, resp.StatusCode)
		return nil, fmt.Errorf("Fail to get page at %s status code %d:", url, resp.StatusCode)
	}
	var resp4Bytes response4Bytes
	err = json.NewDecoder(resp.Body).Decode(&resp4Bytes)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("Fail to decode response at %s: %s\n", url, err.Error())
		return nil, fmt.Errorf("Fail to decode response at %s: %s", url, err.Error())
	}
	var allSignatures []string
	for _, signature := range resp4Bytes.Results {
		//fmt.Printf("%s: %s\n", Event, signature.TextSignature)
		allSignatures = append(allSignatures, signature.TextSignature)
	}
	return allSignatures, nil
}

var indexSignatures = map[string][]string{}

func loadFileEvents() error {
	var allSignatures []signature
	err := json.Unmarshal(embeddedFiles.Get4bytesEvents(), &allSignatures)
	if err != nil {
		//log.Fatalf("Fail to unmarshal: %s", err.Error())
		return fmt.Errorf("Fail to unmarshal: %s", err.Error())
	}
	// Create index for signatures to avoid complexity O(n²) -> O(2n)
	indexSignatures = make(map[string][]string, len(allSignatures))
	for _, signature := range allSignatures {
		indexSignatures[signature.HexSignature] = append(indexSignatures[signature.HexSignature], signature.TextSignature)
	}
	return nil
}

func ReverseEventFile(Event string) ([]string, error) {
	if len(indexSignatures) == 0 {
		err := loadFileEvents()
		if err != nil {
			return nil, err
		}
	}
	var signatures []string
	if textSignature, ok := indexSignatures[Event]; ok {
		signatures = textSignature
		//fmt.Printf("%s: %s\n", arg, strings.Join(textSignature, "\n"+selec+": "))
		// We continue, one Event can match multiple signatures
		//break
	}
	return signatures, nil

}
