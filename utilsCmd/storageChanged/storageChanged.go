package storageChanged

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func DebugTransaction(rpcNode, hash string) ([]byte, error) {
	//json := `{"id": 1, "method": "debug_traceTransaction", "params": ["` + hash + `", {"disableStack": true, "disableMemory": true, "disableStorage": true}]}`
	jsonPayload := `{"id": 1, "method": "debug_traceTransaction", "params": ["` + hash + `"]}`
	jsonByte := []byte(jsonPayload)
	req, err := http.NewRequest("POST", rpcNode, bytes.NewBuffer(jsonByte))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body, err
}
