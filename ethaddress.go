package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EthAddress represents a nullable go-ethereum/common.Address.
type EthAddress struct {
	EthAddress ethcommon.Address
	Valid      bool
}

// NewEthAddress returns a new EthAddress.
func NewEthAddress(address ethcommon.Address, valid bool) EthAddress {
	return EthAddress{
		EthAddress: address,
		Valid:      valid,
	}
}

// NullableString returns the value as a String.
func (n EthAddress) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.EthAddress.String(), true)
}

// Value implements driver.Valuer.
// It returns the driver.Value returned by go-ethereum/common.Address.Value, or nil if invalid.
func (n EthAddress) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.EthAddress.Value()
}

// Scan implements sql.Scanner.
// It accepts any value supported by go-ethereum/common.Address.Scan, or nil.
func (n *EthAddress) Scan(src any) error {
	if src == nil {
		n.EthAddress, n.Valid = ethcommon.Address{}, false

		return nil
	}

	if err := n.EthAddress.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements json.Marshaler.
// It returns the JSON encoding of the string returned by go-ethereum/common.Address.Hex, or null if invalid.
func (n EthAddress) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.EthAddress.Hex())
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts any JSON value supported by go-ethereum/common.Address, or null.
func (n *EthAddress) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.EthAddress, n.Valid = ethcommon.Address{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.EthAddress); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
