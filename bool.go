package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"go.yaml.in/yaml/v3"
)

var (
	_ driver.Valuer        = Bool{}
	_ sql.Scanner          = &Bool{}
	_ json.MarshalerTo     = Bool{}
	_ json.Marshaler       = Bool{}
	_ yaml.Marshaler       = Bool{}
	_ json.UnmarshalerFrom = &Bool{}
	_ json.Unmarshaler     = &Bool{}
	_ yaml.Unmarshaler     = &Bool{}
)

// Bool represents a nullable bool.
type Bool struct {
	sql.NullBool
}

// NewBool returns a new [Bool].
func NewBool(b bool, valid bool) Bool {
	return Bool{
		Bool:  b,
		Valid: valid,
	}
}

// NewBoolFromBoolPtr returns a new [Bool] from a bool pointer.
// It returns an invalid [Bool] if b is nil.
func NewBoolFromBoolPtr(b *bool) Bool {
	if b == nil {
		return NewBool(false, false)
	}

	return NewBool(*b, true)
}

// BoolPtr returns a pointer to a copy of the underlying bool, or nil if n is invalid.
func (n Bool) BoolPtr() *bool {
	if !n.Valid {
		return nil
	}

	return &n.Bool
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as an unquoted boolean string (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Bool) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.Bool)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Bool.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as a bool (or nil if n is invalid).
func (n Bool) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Bool, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes an unquoted boolean string or the unquoted string null from dec into n; the latter makes n invalid.
func (n *Bool) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindTrue, jsontext.KindFalse:
		var b bool
		if err := json.UnmarshalDecode(dec, &b); err != nil {
			return fmt.Errorf("invalid boolean: %w", err)
		}

		n.Bool, n.Valid = b, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Bool, n.Valid = false, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Bool.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Bool) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n as a boolean; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Bool) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Bool, n.Valid = false, false

		return nil
	}

	var b bool
	if err := value.Decode(&b); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	n.Bool, n.Valid = b, true

	return nil
}
