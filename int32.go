package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
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

// NewInt32FromPtr returns a new Int32 from a pointer.
func NewInt32FromPtr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}

	return NewInt32(*i, true)
}

// MarshalJSON implements the json.Marshaler interface.
func (n Int32) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Int32)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
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
