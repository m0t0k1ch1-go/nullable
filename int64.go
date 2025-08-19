package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"github.com/samber/oops"
)

// Int64 is a nullable int64.
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

// NewInt64FromInt64Ptr returns a new Int64 from an int64 pointer.
func NewInt64FromInt64Ptr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}

	return NewInt64(*i, true)
}

// Int64Ptr returns the int64 pointer.
func (n Int64) Int64Ptr() *int64 {
	if !n.Valid {
		return nil
	}

	return &n.Int64
}

// MarshalJSON implements the json.Marshaler interface.
func (n Int64) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(n.Int64)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Int64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Int64, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Int64); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}
