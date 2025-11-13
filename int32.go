package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// Int32 represents a nullable int32 wrapping sql.NullInt32.
type Int32 struct {
	sql.NullInt32
}

// NewInt32 returns a new Int32.
func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		sql.NullInt32{
			Int32: i,
			Valid: valid,
		},
	}
}

// NewInt32FromInt32Ptr returns a new Int32 from a *int32
// It captures the value at call time; a nil pointer is treated as invalid.
func NewInt32FromInt32Ptr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}

	return NewInt32(*i, true)
}

// Int32Ptr returns the value as a *int32, or nil if invalid.
// The pointer refers to a copy.
func (n Int32) Int32Ptr() *int32 {
	if !n.Valid {
		return nil
	}

	return &n.Int32
}

// MarshalJSON implements json.Marshaler.
// It returns the value as a JSON number, or null if invalid.
func (n Int32) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Int32)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON number or null.
func (n *Int32) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Int32, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Int32); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
