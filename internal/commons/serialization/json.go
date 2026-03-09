// Package serialization provides serialization utilities.
package serialization

import "encoding/json"

// ToJSON serializes v to a JSON byte slice.
func ToJSON(v any) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON deserializes JSON bytes into v.
func FromJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// ToJSONString serializes v to a JSON string.
func ToJSONString(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
