package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// Bool is a nullable bool.
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

// NewBoolFromPtr returns a new Bool from a pointer.
func NewBoolFromPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}

	return NewBool(*b, true)
}

// MarshalJSON implements the json.Marshaler interface.
func (n Bool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Bool)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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
