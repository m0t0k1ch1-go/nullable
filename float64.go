package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// Float64 represents a nullable float64 wrapping sql.NullFloat64.
type Float64 struct {
	sql.NullFloat64
}

// NewFloat64 returns a new Float64.
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		sql.NullFloat64{
			Float64: f,
			Valid:   valid,
		},
	}
}

// NewFloat64FromFloat64Ptr returns a new Float64 from a *float64
// It captures the value at call time; a nil pointer is treated as invalid.
func NewFloat64FromFloat64Ptr(f *float64) Float64 {
	if f == nil {
		return NewFloat64(0, false)
	}

	return NewFloat64(*f, true)
}

// Float64Ptr returns the value as a *float64, or nil if invalid.
// The pointer refers to a copy.
func (n Float64) Float64Ptr() *float64 {
	if !n.Valid {
		return nil
	}

	return &n.Float64
}

// MarshalJSON implements json.Marshaler.
// It returns the value as a JSON number, or null if invalid.
func (n Float64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Float64)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON number or null.
func (n *Float64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Float64, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Float64); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
