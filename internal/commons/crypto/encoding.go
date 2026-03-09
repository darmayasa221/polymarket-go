package crypto

import "encoding/base64"

// Base64Encode encodes bytes to base64 URL-safe string.
func Base64Encode(data []byte) string {
	return base64.URLEncoding.EncodeToString(data)
}

// Base64Decode decodes a base64 URL-safe string to bytes.
func Base64Decode(s string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(s)
}
