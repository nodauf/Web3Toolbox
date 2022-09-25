package embeddedFiles

import (
	_ "embed"
)

//go:embed 4bytesSignatures.json
var signaturesContent []byte

//go:embed 4bytesEvents.json
var eventsContent []byte

// Get4bytesSignatures returns the content of the embedded file 4bytesSignatures.json
func Get4bytesSignatures() []byte {
	return signaturesContent
}

// Get4bytesEvents returns the content of the embedded file 4bytesEvents.json
func Get4bytesEvents() []byte {
	return eventsContent
}
