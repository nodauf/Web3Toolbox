package embeddedFiles

import (
	_ "embed"
)

//go:embed 4bytesSignatures.json
var signaturesContent []byte

// Get4bytesSignatures returns the content of the embedded file 4bytesSignatures.json
func Get4bytesSignatures() []byte {
	return signaturesContent
}
