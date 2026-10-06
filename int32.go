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
	_ driver.Valuer        = Int32{}
	_ sql.Scanner          = &Int32{}
	_ json.MarshalerTo     = Int32{}
	_ json.Marshaler       = Int32{}
	_ yaml.Marshaler       = Int32{}
	_ json.UnmarshalerFrom = &Int32{}
	_ json.Unmarshaler     = &Int32{}
	_ yaml.Unmarshaler     = &Int32{}
)

// Int32 represents a nullable int32.
type Int32 struct {
	sql.NullInt32
}

// NewInt32 returns a new [Int32].
func NewInt32(i int32, valid bool) Int32 {
	return Int32{
		Int32: i,
		Valid: valid,
	}
}

// NewInt32FromInt32Ptr returns a new [Int32] from an int32 pointer.
// It returns an invalid [Int32] if i is nil.
func NewInt32FromInt32Ptr(i *int32) Int32 {
	if i == nil {
		return NewInt32(0, false)
	}

	return NewInt32(*i, true)
}

// Int32Ptr returns a pointer to a copy of the underlying int32, or nil if n is invalid.
func (n Int32) Int32Ptr() *int32 {
	if !n.Valid {
		return nil
	}

	return &n.Int32
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as an unquoted decimal string (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Int32) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.Int32)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Int32.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Int32) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as an int32 (or nil if n is invalid).
func (n Int32) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Int32, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes an unquoted decimal string or the unquoted string null from dec into n; the latter makes n invalid.
func (n *Int32) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindNumber:
		var i int32
		if err := json.UnmarshalDecode(dec, &i); err != nil {
			return fmt.Errorf("invalid number: %w", err)
		}

		n.Int32, n.Valid = i, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Int32, n.Valid = 0, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Int32.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Int32) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n as an int32; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Int32) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Int32, n.Valid = 0, false

		return nil
	}

	var i int32
	if err := value.Decode(&i); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	n.Int32, n.Valid = i, true

	return nil
}
