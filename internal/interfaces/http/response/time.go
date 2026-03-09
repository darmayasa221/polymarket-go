package response

import "time"

// JSONTime wraps time.Time and serializes to RFC3339 (no nanoseconds) in JSON responses.
type JSONTime time.Time

// MarshalJSON formats the time as RFC3339 without nanoseconds.
func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).UTC().Format(time.RFC3339)
	return []byte(`"` + formatted + `"`), nil
}

// UnmarshalJSON parses an RFC3339 timestamp into JSONTime.
func (t *JSONTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, string(data[1:len(data)-1]))
	if err != nil {
		return err
	}
	*t = JSONTime(parsed)
	return nil
}
