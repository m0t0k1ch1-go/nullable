package nullable

import (
	"bytes"
	"math/big"
)

// BigInt is a nullable math/big.Int.
type BigInt struct {
	BigInt *big.Int
	Valid  bool
}

// NewBigInt returns a new BigInt.
func NewBigInt(x *big.Int, valid bool) BigInt {
	return BigInt{
		BigInt: x,
		Valid:  valid,
	}
}

// NullableString returns the String.
func (n BigInt) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.BigInt.String(), true)
}

// MarshalJSON implements the json.Marshaler interface.
func (n BigInt) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return n.BigInt.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *BigInt) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.BigInt, n.Valid = nil, false

		return nil
	}

	var x big.Int
	if err := x.UnmarshalJSON(b); err != nil {
		return err
	}

	n.BigInt, n.Valid = &x, true

	return nil
}
