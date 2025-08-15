package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/samber/oops"
)

// EthAddress is a nullable github.com/ethereum/go-ethereum/common.Address.
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

// NullableString returns the String.
func (n EthAddress) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.EthAddress.String(), true)
}

// Value implements the driver.Valuer interface.
func (n EthAddress) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	v, err := n.EthAddress.Value()
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return v, nil
}

// Scan implements the sql.Scanner interface.
func (n *EthAddress) Scan(src any) error {
	if src == nil {
		n.EthAddress, n.Valid = ethcommon.Address{}, false

		return nil
	}

	if err := n.EthAddress.Scan(src); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n EthAddress) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(n.EthAddress.Hex())
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *EthAddress) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.EthAddress, n.Valid = ethcommon.Address{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.EthAddress); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}
