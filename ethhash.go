package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EthHash represents a nullable go-ethereum/common.Hash
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

// NullableString returns the value as a String.
func (n EthHash) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.EthHash.String(), true)
}

// Value implements driver.Valuer.
// It returns the driver.Value returned by go-ethereum/common.Hash.Value, or nil if invalid.
func (n EthHash) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.EthHash.Value()
}

// Scan implements sql.Scanner.
// It accepts any value supported by go-ethereum/common.Hash.Scan, or nil.
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

// MarshalJSON implements json.Marshaler.
// It returns the JSON encoding of the string returned by go-ethereum/common.Hash.Hex, or null if invalid.
func (n EthHash) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.EthHash.Hex())
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts any JSON value supported by go-ethereum/common.Hash, or null.
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
