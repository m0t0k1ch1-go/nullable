package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EthHash is a nullable github.com/ethereum/go-ethereum/common.Hash.
type EthHash struct {
	EthHash ethcommon.Hash
	Valid   bool
}

// NewEthHash returns a new EthHash.
func NewEthHash(h ethcommon.Hash, valid bool) EthHash {
	return EthHash{
		EthHash: h,
		Valid:   valid,
	}
}

// NullableString returns the String.
func (n EthHash) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.EthHash.String(), true)
}

// Value implements the driver.Valuer interface.
func (n EthHash) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.EthHash.Value()
}

// Scan implements the sql.Scanner interface.
func (n *EthHash) Scan(src any) error {
	if src == nil {
		n.EthHash, n.Valid = ethcommon.Hash{}, false

		return nil
	}

	if err := n.EthHash.Scan(src); err != nil {
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

	return json.Marshal(n.EthHash.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *EthHash) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.EthHash, n.Valid = ethcommon.Hash{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.EthHash); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
