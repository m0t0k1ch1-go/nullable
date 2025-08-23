package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"go.yaml.in/yaml/v3"
)

// String represents a nullable string wrapping sql.NullString.
type String struct {
	sql.NullString
}

// NewString returns a new String.
func NewString(s string, valid bool) String {
	return String{
		sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// NewStringFromStringPtr returns a new String from a *string.
// It captures the value at call time; a nil pointer is treated as invalid.
func NewStringFromStringPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}

	return NewString(*s, true)
}

// StringPtr returns the value as a *string, or nil if invalid.
// The pointer refers to a copy.
func (n String) StringPtr() *string {
	if !n.Valid {
		return nil
	}

	return &n.String
}

// MarshalJSON implements json.Marshaler.
// It returns the value as a JSON string, or null if invalid.
func (n String) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.String)
}

// MarshalYAML implements yaml.Marshaler.
// It returns the value as a string, or nil if invalid.
func (n String) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.String, nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON string or null.
func (n *String) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.String, n.Valid = "", false

		return nil
	}

	if err := json.Unmarshal(b, &n.String); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
// It accepts a YAML string or null.
// Note: yaml.v3 may bypass this method for null; handle the explicit !!null tag defensively.
func (n *String) UnmarshalYAML(value *yaml.Node) error {
	if value.Tag == "!!null" {
		n.String, n.Valid = "", false

		return nil
	}

	if err := value.Decode(&n.String); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
