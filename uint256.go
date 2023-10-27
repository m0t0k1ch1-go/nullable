package nullable

import (
	"bytes"
	"database/sql/driver"

	"github.com/m0t0k1ch1-go/bigutil/v2"
)

// Uint256 is a nullable github.com/m0t0k1ch1-go/bigutil.Uint256.
type Uint256 struct {
	Uint256 bigutil.Uint256
	Valid   bool
}

// NewUint256 returns a new Uint256.
func NewUint256(i bigutil.Uint256, valid bool) Uint256 {
	return Uint256{
		Uint256: i,
		Valid:   valid,
	}
}

// NullableString returns the String.
func (i Uint256) NullableString() String {
	if !i.Valid {
		return NewString("", false)
	}

	return NewString(i.Uint256.String(), true)
}

// Value implements the driver.Valuer interface.
func (i Uint256) Value() (driver.Value, error) {
	if !i.Valid {
		return nil, nil
	}

	return i.Uint256.Value()
}

// Scan implements the sql.Scanner interface.
func (i *Uint256) Scan(src any) error {
	if src == nil {
		i.Uint256, i.Valid = bigutil.Uint256{}, false

		return nil
	}

	if err := i.Uint256.Scan(src); err != nil {
		return err
	}

	i.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (i Uint256) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}

	return i.Uint256.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (i *Uint256) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		i.Uint256, i.Valid = bigutil.Uint256{}, false

		return nil
	}

	if err := i.Uint256.UnmarshalJSON(b); err != nil {
		return err
	}

	i.Valid = true

	return nil
}
