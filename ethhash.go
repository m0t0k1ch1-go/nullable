package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// EthHash is a nullable ethcommon.Hash.
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

// NullableString returns the String representation of the EthHash.
func (h EthHash) NullableString() String {
	if !h.Valid {
		return NewString("", false)
	}

	return NewString(h.Hash.String(), true)
}

// Value implements the driver.Valuer interface.
func (h EthHash) Value() (driver.Value, error) {
	if !h.Valid {
		return nil, nil
	}

	return h.Hash.Value()
}

// Scan implements the sql.Scanner interface.
func (h *EthHash) Scan(src any) error {
	if src == nil {
		h.Hash, h.Valid = ethcommon.Hash{}, false

		return nil
	}

	if err := h.Hash.Scan(src); err != nil {
		return err
	}

	h.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (h EthHash) MarshalJSON() ([]byte, error) {
	if !h.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(h.Hash.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (h *EthHash) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		h.Hash, h.Valid = ethcommon.Hash{}, false

		return nil
	}

	if err := h.Hash.UnmarshalJSON(b); err != nil {
		return err
	}

	h.Valid = true

	return nil
}
