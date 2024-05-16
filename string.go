package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

// NewStringFromPtr returns a new String from a pointer.
func NewStringFromPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}

	return NewString(*s, true)
}

// MarshalJSON implements the json.Marshaler interface.
func (n String) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.String)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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
		return err
	}

	n.Valid = true

	return nil
}
