package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EthHash is a nullable github.com/ethereum/go-ethereum/common.Hash.
type EthHash struct {
	Hash  ethcommon.Hash
	Valid bool
}

// NewEthHash returns a new EthHash.
func NewEthHash(h ethcommon.Hash, valid bool) EthHash {
	return EthHash{
		Hash:  h,
		Valid: valid,
	}
}

// NullableString returns the String.
func (n EthHash) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Hash.String(), true)
}

// Value implements the driver.Valuer interface.
func (n EthHash) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Hash.Value()
}

// Scan implements the sql.Scanner interface.
func (n *EthHash) Scan(src any) error {
	if src == nil {
		n.Hash, n.Valid = ethcommon.Hash{}, false

		return nil
	}

	if err := n.Hash.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n EthHash) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Hash.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *EthHash) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Hash, n.Valid = ethcommon.Hash{}, false

		return nil
	}

	if err := n.Hash.UnmarshalJSON(b); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
