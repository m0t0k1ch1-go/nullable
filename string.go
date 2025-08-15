package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"github.com/samber/oops"
)

// String is a nullable string.
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
func NewStringFromStringPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}

	return NewString(*s, true)
}

// StringPtr returns the string pointer.
func (n String) StringPtr() *string {
	if !n.Valid {
		return nil
	}

	return &n.String
}

// MarshalJSON implements the json.Marshaler interface.
func (n String) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(n.String)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *String) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.String, n.Valid = "", false

		return nil
	}

	if err := json.Unmarshal(b, &n.String); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}

// MarshalYAML implements the yaml.Marshaler interface.
func (n String) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.String, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (n *String) UnmarshalYAML(unmarshal func(any) error) error {
	if err := unmarshal(&n.String); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}
