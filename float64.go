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
	_ driver.Valuer        = Float64{}
	_ sql.Scanner          = &Float64{}
	_ json.MarshalerTo     = Float64{}
	_ json.Marshaler       = Float64{}
	_ yaml.Marshaler       = Float64{}
	_ json.UnmarshalerFrom = &Float64{}
	_ json.Unmarshaler     = &Float64{}
	_ yaml.Unmarshaler     = &Float64{}
)

// Float64 represents a nullable float64.
type Float64 struct {
	sql.NullFloat64
}

// NewFloat64 returns a new [Float64].
func NewFloat64(f float64, valid bool) Float64 {
	return Float64{
		Float64: f,
		Valid:   valid,
	}
}

// NewFloat64FromFloat64Ptr returns a new [Float64] from a float64 pointer.
// It returns an invalid [Float64] if f is nil.
func NewFloat64FromFloat64Ptr(f *float64) Float64 {
	if f == nil {
		return NewFloat64(0, false)
	}

	return NewFloat64(*f, true)
}

// Float64Ptr returns a pointer to a copy of the underlying float64, or nil if n is invalid.
func (n Float64) Float64Ptr() *float64 {
	if !n.Valid {
		return nil
	}

	return &n.Float64
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as an unquoted decimal string (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Float64) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.Float64)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Float64.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as a float64 (or nil if n is invalid).
func (n Float64) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Float64, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes an unquoted decimal string or the unquoted string null from dec into n; the latter makes n invalid.
func (n *Float64) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindNumber:
		var f float64
		if err := json.UnmarshalDecode(dec, &f); err != nil {
			return fmt.Errorf("invalid float64: %w", err)
		}

		n.Float64, n.Valid = f, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Float64, n.Valid = 0, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Float64.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Float64) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n as a float64; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Float64) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Float64, n.Valid = 0, false

		return nil
	}

	var f float64
	if err := value.Decode(&f); err != nil {
		return fmt.Errorf("invalid float64: %w", err)
	}

	n.Float64, n.Valid = f, true

	return nil
}
