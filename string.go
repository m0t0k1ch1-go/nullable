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
	_ driver.Valuer        = String{}
	_ sql.Scanner          = &String{}
	_ json.MarshalerTo     = String{}
	_ json.Marshaler       = String{}
	_ yaml.Marshaler       = String{}
	_ json.UnmarshalerFrom = &String{}
	_ json.Unmarshaler     = &String{}
	_ yaml.Unmarshaler     = &String{}
)

// String represents a nullable string.
type String struct {
	sql.NullString
}

// NewString returns a new [String].
func NewString(s string, valid bool) String {
	return String{
		String: s,
		Valid:  valid,
	}
}

// NewStringFromStringPtr returns a new [String] from a string pointer.
// It returns an invalid [String] if s is nil.
func NewStringFromStringPtr(s *string) String {
	if s == nil {
		return NewString("", false)
	}

	return NewString(*s, true)
}

// StringPtr returns a pointer to a copy of the underlying string, or nil if n is invalid.
func (n String) StringPtr() *string {
	if !n.Valid {
		return nil
	}

	return &n.String
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as a quoted string (or null if n is invalid) and writes it to enc.
func (n String) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.String)
}

// MarshalJSON implements [json.Marshaler].
// It is like [String.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n String) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as a string (or nil if n is invalid).
func (n String) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.String, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a quoted string or null from dec into n; null makes n invalid.
func (n *String) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindString:
		var s string
		if err := json.UnmarshalDecode(dec, &s); err != nil {
			return fmt.Errorf("invalid string: %w", err)
		}

		n.String, n.Valid = s, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.String, n.Valid = "", false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [String.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *String) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n as a string; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *String) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.String, n.Valid = "", false

		return nil
	}

	var s string
	if err := value.Decode(&s); err != nil {
		return fmt.Errorf("invalid string: %w", err)
	}

	n.String, n.Valid = s, true

	return nil
}
