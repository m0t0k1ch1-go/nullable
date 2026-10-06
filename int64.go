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
	_ driver.Valuer        = Int64{}
	_ sql.Scanner          = &Int64{}
	_ json.MarshalerTo     = Int64{}
	_ json.Marshaler       = Int64{}
	_ yaml.Marshaler       = Int64{}
	_ json.UnmarshalerFrom = &Int64{}
	_ json.Unmarshaler     = &Int64{}
	_ yaml.Unmarshaler     = &Int64{}
)

// Int64 represents a nullable int64.
type Int64 struct {
	sql.NullInt64
}

// NewInt64 returns a new [Int64].
func NewInt64(i int64, valid bool) Int64 {
	return Int64{
		Int64: i,
		Valid: valid,
	}
}

// NewInt64FromInt64Ptr returns a new [Int64] from an int64 pointer.
// It returns an invalid [Int64] if i is nil.
func NewInt64FromInt64Ptr(i *int64) Int64 {
	if i == nil {
		return NewInt64(0, false)
	}

	return NewInt64(*i, true)
}

// Int64Ptr returns a pointer to a copy of the underlying int64, or nil if n is invalid.
func (n Int64) Int64Ptr() *int64 {
	if !n.Valid {
		return nil
	}

	return &n.Int64
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as an unquoted decimal string (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Int64) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.Int64)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Int64.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as an int64 (or nil if n is invalid).
func (n Int64) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Int64, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes an unquoted decimal string or the unquoted string null from dec into n; the latter makes n invalid.
func (n *Int64) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindNumber:
		var i int64
		if err := json.UnmarshalDecode(dec, &i); err != nil {
			return fmt.Errorf("invalid number: %w", err)
		}

		n.Int64, n.Valid = i, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Int64, n.Valid = 0, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Int64.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Int64) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n as an int64; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Int64) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Int64, n.Valid = 0, false

		return nil
	}

	var i int64
	if err := value.Decode(&i); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	n.Int64, n.Valid = i, true

	return nil
}
