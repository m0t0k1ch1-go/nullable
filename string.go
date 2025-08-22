package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

// NewStringFromStringPtr returns a new String from a string pointer.
// A nil pointer is treated as invalid.
func NewStringFromStringPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}

	return NewString(*s, true)
}

// StringPtr returns the value as a string pointer, or nil if invalid.
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

// MarshalYAML implements yaml.Marshaler.
// It returns the value as a string, or nil if invalid.
func (n String) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.String, nil
}

// UnmarshalYAML implements yaml.Unmarshaler for yaml.v2.
// It accepts a YAML string.
// Deprecated: Kept for yaml.v2 compatibility.
func (n *String) UnmarshalYAML(unmarshal func(any) error) error {
	if err := unmarshal(&n.String); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
