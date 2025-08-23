package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// Uint64 represents a nullable uint64.
type Uint64 struct {
	Uint64 uint64
	Valid  bool
}

// NewUint64 returns a new Uint64.
func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		Uint64: i,
		Valid:  valid,
	}
}

// NewUint64FromUint64Ptr returns a new Uint64 from a *uint64.
// It captures the value at call time; a nil pointer is treated as invalid.
func NewUint64FromUint64Ptr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}

	return NewUint64(*i, true)
}

// Uint64Ptr returns the value as a *uint64, or nil if invalid.
// The pointer refers to a copy.
func (n Uint64) Uint64Ptr() *uint64 {
	if !n.Valid {
		return nil
	}

	return &n.Uint64
}

// Value implements driver.Valuer.
// It returns the value as a uint64, or nil if invalid.
func (n Uint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint64, nil
}

// Scan implements sql.Scanner.
// It accepts one of the following:
//   - int64 (non-negative)
//   - uint64
//   - []byte (non-negative decimal string)
//   - nil
func (n *Uint64) Scan(src any) error {
	if src == nil {
		n.Uint64, n.Valid = 0, false

		return nil
	}

	switch v := src.(type) {

	case int64:
		if v < 0 {
			return errors.New("invalid source: negative int64")
		}

		n.Uint64, n.Valid = uint64(v), true

		return nil

	case uint64:
		n.Uint64, n.Valid = v, true

		return nil

	case []byte:
		if len(v) == 0 {
			return errors.New("invalid source: empty []byte")
		}

		i, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid source: %w", err)
		}

		n.Uint64, n.Valid = i, true

		return nil

	default:
		return fmt.Errorf("unsupported source type: %T", src)
	}
}

// MarshalJSON implements json.Marshaler.
// It returns the value as a JSON number, or null if invalid.
func (n Uint64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Uint64)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts a JSON number (non-negative integer) or null.
func (n *Uint64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Uint64, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Uint64); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
