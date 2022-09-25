package reverseEvent

type signature struct {
	ID            int    `json:"id"`
	TextSignature string `json:"text_signature"`
	HexSignature  string `json:"hex_signature"`
}

type response4Bytes struct {
	Results []signature `json:"results"`
}
