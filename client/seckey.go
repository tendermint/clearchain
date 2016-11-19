package client

import (
	"encoding/base64"
	"fmt"
)

// Encode converts base64 encoded binary block into string.
func Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Decode converts string data back into a base64 encoded binary data block.
func Decode(base64EncodedString string) []byte {
	base64DecodedString, err := base64.StdEncoding.DecodeString(base64EncodedString)

	if err != nil {
		panic(fmt.Sprintf("Error (%v) during decoding data: %v", base64EncodedString, err))
	}

	return base64DecodedString
}
