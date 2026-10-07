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
	_ driver.Valuer        = EthAddress{}
	_ sql.Scanner          = &EthAddress{}
	_ json.MarshalerTo     = EthAddress{}
	_ json.Marshaler       = EthAddress{}
	_ yaml.Marshaler       = EthAddress{}
	_ json.UnmarshalerFrom = &EthAddress{}
	_ json.Unmarshaler     = &EthAddress{}
	_ yaml.Unmarshaler     = &EthAddress{}
)

// EthAddress represents a nullable [ethcommon.Address].
type EthAddress struct {
	EthAddress ethcommon.Address
	Valid      bool
}

// NewEthAddress returns a new [EthAddress].
func NewEthAddress(a ethcommon.Address, valid bool) EthAddress {
	return EthAddress{
		EthAddress: a,
		Valid:      valid,
	}
}

// NullableString returns the string returned by [ethcommon.Address.Hex] as a [String], or an invalid [String] if n is invalid.
func (n EthAddress) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.EthAddress.Hex(), true)
}

// Value implements [driver.Valuer].
// It encodes n by delegating to [ethcommon.Address.Value] (or as nil if n is invalid).
func (n EthAddress) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.EthAddress.Value()
}

// Scan implements [sql.Scanner].
// It decodes src into n by delegating to [ethcommon.Address.Scan]; nil makes n invalid.
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

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as the quoted string returned by [ethcommon.Address.Hex] (or as the unquoted string null if n is invalid) and writes it to enc.
func (n EthAddress) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.EthAddress.Hex())
}

// MarshalJSON implements [json.Marshaler].
// It is like [EthAddress.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n EthAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as the quoted string returned by [ethcommon.Address.Hex] (or nil if n is invalid).
func (n EthAddress) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Style: yaml.DoubleQuotedStyle,
		Value: n.EthAddress.Hex(),
	}, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a quoted string from dec into n by delegating to [ethcommon.Address.UnmarshalText]; the unquoted string null makes n invalid.
func (n *EthAddress) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindString:
		var s string
		if err := json.UnmarshalDecode(dec, &s); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		var a ethcommon.Address
		if err := a.UnmarshalText([]byte(s)); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		n.EthAddress, n.Valid = a, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.EthAddress, n.Valid = ethcommon.Address{}, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [EthAddress.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *EthAddress) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n by delegating to [ethcommon.Address.UnmarshalText]; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *EthAddress) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.EthAddress, n.Valid = ethcommon.Address{}, false

		return nil
	}

	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	var a ethcommon.Address
	if err := a.UnmarshalText([]byte(s)); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	n.EthAddress, n.Valid = a, true

	return nil
}
