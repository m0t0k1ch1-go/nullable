package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// Bool represents a nullable bool wrapping sql.NullBool.
type Bool struct {
	sql.NullBool
}

// NewBool returns a new Bool.
func NewBool(b bool, valid bool) Bool {
	return Bool{
		sql.NullBool{
			Bool:  b,
			Valid: valid,
		},
	}
}

// NewBoolFromBoolPtr returns a new Bool from a bool pointer.
// A nil pointer is treated as invalid.
func NewBoolFromBoolPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}

	return NewBool(*b, true)
}

// BoolPtr returns the value as a bool pointer, or nil if invalid.
// The pointer refers to a copy.
func (n Bool) BoolPtr() *bool {
	if !n.Valid {
		return nil
	}

	return &n.Bool
}

// MarshalJSON implements json.Marshaler.
// It returns the value as a JSON boolean, or null if invalid.
func (n Bool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Bool)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON boolean or null.
func (n *Bool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Bool, n.Valid = false, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Bool); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
