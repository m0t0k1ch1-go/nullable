package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/samber/oops"
)

// Uint64 is a nullable uint64.
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

// NewUint64FromUint64Ptr returns a new Uint64 from a uint64 pointer.
func NewUint64FromUint64Ptr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}

	return NewUint64(*i, true)
}

// Uint64Ptr returns the uint64 pointer.
func (n Uint64) Uint64Ptr() *uint64 {
	if !n.Valid {
		return nil
	}

	return &n.Uint64
}

// Value implements the driver.Valuer interface.
func (n Uint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint64, nil
}

// Scan implements the sql.Scanner interface.
func (n *Uint64) Scan(src any) error {
	if src == nil {
		n.Uint64, n.Valid = 0, false

		return nil
	}

	switch v := src.(type) {

	case int64:
		if v < 0 {
			return oops.New("src must not be negative")
		}

		n.Uint64 = uint64(v)

	case uint64:
		n.Uint64 = v

	case []byte:
		i, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return oops.Wrapf(err, "failed to parse src as uint64")
		}

		n.Uint64 = i

	default:
		return oops.Errorf("unexpected src type: %T", src)
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n Uint64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(n.Uint64)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Uint64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Uint64, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Uint64); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}
