package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// Int64 represents a nullable int64 wrapping sql.NullInt64.
type Int64 struct {
	sql.NullInt64
}

// NewInt64 returns a new Int64.
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		sql.NullInt64{
			Int64: i,
			Valid: valid,
		},
	}
}

// NewInt64FromInt64Ptr returns a new Int64 from a *int64
// It captures the value at call time; a nil pointer is treated as invalid.
func NewInt64FromInt64Ptr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}

	return NewInt64(*i, true)
}

// Int64Ptr returns the value as a *int64, or nil if invalid.
// The pointer refers to a copy.
func (n Int64) Int64Ptr() *int64 {
	if !n.Valid {
		return nil
	}

	return &n.Int64
}

// MarshalJSON implements json.Marshaler.
// It returns the value as a JSON number, or null if invalid.
func (n Int64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Int64)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON number or null.
func (n *Int64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Int64, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Int64); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
