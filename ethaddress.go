package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EthAddress is a nullable github.com/ethereum/go-ethereum/common.Address.
type EthAddress struct {
	Address ethcommon.Address
	Valid   bool
}

// NewEthAddress returns a new EthAddress.
func NewEthAddress(address ethcommon.Address, valid bool) EthAddress {
	return EthAddress{
		Address: address,
		Valid:   valid,
	}
}

// NullableString returns the String.
func (n EthAddress) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Address.String(), true)
}

// Value implements the driver.Valuer interface.
func (n EthAddress) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Address.Value()
}

// Scan implements the sql.Scanner interface.
func (n *EthAddress) Scan(src any) error {
	if src == nil {
		n.Address, n.Valid = ethcommon.Address{}, false

		return nil
	}

	if err := n.Address.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n EthAddress) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Address.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *EthAddress) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Address, n.Valid = ethcommon.Address{}, false

		return nil
	}

	if err := n.Address.UnmarshalJSON(b); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
