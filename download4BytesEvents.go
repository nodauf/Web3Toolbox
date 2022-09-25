package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

type signature struct {
	ID            int    `json:"id"`
	TextSignature string `json:"text_signature"`
	HexSignature  string `json:"hex_signature"`
}

type response4Bytes struct {
	Results []signature `json:"results"`
}

func main() {
	var thread, pageMax int
	flag.IntVar(&thread, "thread", 10, "Number of threads")
	flag.IntVar(&pageMax, "pageMax", 8095, "The last page to retrieve")
	flag.Parse()

	var bar = progressbar.Default(int64(pageMax))
	var allSignatures []signature
	pages := make(chan int)
	results := make(chan []signature)
	var mux sync.Mutex
	var wg sync.WaitGroup
	// start workers pool with  threads
	for i := 0; i < thread; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			go download4BytesPage(pages, results, bar)

			result := <-results
			mux.Lock()
			allSignatures = append(allSignatures, result...)
			mux.Unlock()
		}()

	}
	//8095
	for page := 1; page <= pageMax; page++ {
		pages <- page

	}
	// Timeout to close chan only when all page are processed
pageChecker:
	for {
		select {
		case page := <-pages:
			// There are more pages to process. Put back the page into the chan and wait
			pages <- page
		default:
			// All pages have been processed. Exit the loop for
			break pageChecker
		}
		time.Sleep(30 * time.Second)
	}
	close(pages)
	wg.Wait()

	signatureJSON, _ := json.MarshalIndent(allSignatures, "", "  ")
	err := ioutil.WriteFile("4bytesSignatures.json", signatureJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func download4BytesPage(pages chan int, results chan<- []signature, bar *progressbar.ProgressBar) {
	var baseURL = "https://www.4byte.directory/api/v1/signatures/?page=%d"

	var allSignatures []signature
	defer func() { results <- allSignatures }()
	var resp4Bytes response4Bytes
	for page := range pages {

		url := fmt.Sprintf(baseURL, page)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Fail to get page at %s: %s\n", url, err.Error())
			// We need to process the page again
			pages <- page
			continue
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Fail to get page at %s status code %d:\n", url, resp.StatusCode)
			resp.Body.Close()
			// We need to process the page again
			pages <- page
			// If the server is overload, wait a bit before next request
			if resp.StatusCode == 502 {
				time.Sleep(10 * time.Second)
			}
			continue
		}
		//a, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println(string(a))
		err = json.NewDecoder(resp.Body).Decode(&resp4Bytes)
		//err = json.Unmarshal(a, &resp4Bytes)
		if err != nil {
			fmt.Printf("Fail to decode response at %s: %s\n", url, err.Error())
			// We need to process the page again
			pages <- page
			continue
		}

		// We reach the end
		if len(resp4Bytes.Results) == 0 {
			fmt.Println("Reach the end")
		}
		allSignatures = append(allSignatures, resp4Bytes.Results...)
		bar.Add(1)
		//allSignatures = append(allSignatures, resp4Bytes.Results...)
		//result <- resp4Bytes

	}

}
