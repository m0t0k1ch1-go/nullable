package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"go.yaml.in/yaml/v3"
)

var (
	_ driver.Valuer        = EthHash{}
	_ sql.Scanner          = &EthHash{}
	_ json.MarshalerTo     = EthHash{}
	_ json.Marshaler       = EthHash{}
	_ yaml.Marshaler       = EthHash{}
	_ json.UnmarshalerFrom = &EthHash{}
	_ json.Unmarshaler     = &EthHash{}
	_ yaml.Unmarshaler     = &EthHash{}
)

// EthHash represents a nullable [ethcommon.Hash].
type EthHash struct {
	EthHash ethcommon.Hash
	Valid   bool
}

// NewEthHash returns a new [EthHash].
func NewEthHash(h ethcommon.Hash, valid bool) EthHash {
	return EthHash{
		EthHash: h,
		Valid:   valid,
	}
}

// NullableString returns the string returned by [ethcommon.Hash.Hex] as a [String], or an invalid [String] if n is invalid.
func (n EthHash) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.EthHash.Hex(), true)
}

// Value implements [driver.Valuer].
// It encodes n by delegating to [ethcommon.Hash.Value] (or as nil if n is invalid).
func (n EthHash) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.EthHash.Value()
}

// Scan implements [sql.Scanner].
// It decodes src into n by delegating to [ethcommon.Hash.Scan]; nil makes n invalid.
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

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as the quoted string returned by [ethcommon.Hash.Hex] (or as the unquoted string null if n is invalid) and writes it to enc.
func (n EthHash) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.EthHash.Hex())
}

// MarshalJSON implements [json.Marshaler].
// It is like [EthHash.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n EthHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as the quoted string returned by [ethcommon.Hash.Hex] (or nil if n is invalid).
func (n EthHash) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.DoubleQuotedStyle,
		Value: n.EthHash.Hex(),
	}, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a quoted string from dec into n by delegating to [ethcommon.Hash.UnmarshalText]; the unquoted string null makes n invalid.
func (n *EthHash) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindString:
		var s string
		if err := json.UnmarshalDecode(dec, &s); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		var h ethcommon.Hash
		if err := h.UnmarshalText([]byte(s)); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		n.EthHash, n.Valid = h, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.EthHash, n.Valid = ethcommon.Hash{}, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [EthHash.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *EthHash) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n by delegating to [ethcommon.Hash.UnmarshalText]; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *EthHash) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.EthHash, n.Valid = ethcommon.Hash{}, false

		return nil
	}

	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	var h ethcommon.Hash
	if err := h.UnmarshalText([]byte(s)); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	n.EthHash, n.Valid = h, true

	return nil
}
