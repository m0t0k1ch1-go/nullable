package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/m0t0k1ch1-go/bigutil/v3"
)

// Uint256 represents a nullable bigutil.Uint256.
type Uint256 struct {
	Uint256 bigutil.Uint256
	Valid   bool
}

// NewUint256 returns a new Uint256.
func NewUint256(x256 bigutil.Uint256, valid bool) Uint256 {
	return Uint256{
		Uint256: x256,
		Valid:   valid,
	}
}

// NullableString returns the value as a String.
func (n Uint256) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Uint256.String(), true)
}

// Value implements driver.Valuer.
// It returns the driver.Value representation of bigutil.Uint256, or nil if invalid.
func (n Uint256) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint256.Value()
}

// Scan implements sql.Scanner.
// It accepts a value supported by bigutil.Uint256.Scan, or nil.
func (n *Uint256) Scan(src any) error {
	if src == nil {
		n.Uint256, n.Valid = bigutil.Uint256{}, false

		return nil
	}

	if err := n.Uint256.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements json.Marshaler.
// It returns the value as the JSON encoding of bigutil.Uint256, or null if invalid.
func (n Uint256) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Uint256)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts the JSON encoding of bigutil.Uint256, or null.
func (n *Uint256) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Uint256, n.Valid = bigutil.Uint256{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Uint256); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
