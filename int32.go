package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"github.com/samber/oops"
)

// Int32 is a nullable int32.
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

// NewInt32FromInt32Ptr returns a new Int32 from a int32 pointer.
func NewInt32FromInt32Ptr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}

	return NewInt32(*i, true)
}

// Int32Ptr returns the int32 pointer.
func (n Int32) Int32Ptr() *int32 {
	if !n.Valid {
		return nil
	}

	return &n.Int32
}

// MarshalJSON implements the json.Marshaler interface.
func (n Int32) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(n.Int32)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Int32) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Int32, n.Valid = 0, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Int32); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}
