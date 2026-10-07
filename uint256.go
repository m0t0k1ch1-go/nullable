package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/m0t0k1ch1-go/bigutil/v3"
	"go.yaml.in/yaml/v3"
)

var (
	_ driver.Valuer        = Uint256{}
	_ sql.Scanner          = &Uint256{}
	_ json.MarshalerTo     = Uint256{}
	_ json.Marshaler       = Uint256{}
	_ yaml.Marshaler       = Uint256{}
	_ json.UnmarshalerFrom = &Uint256{}
	_ json.Unmarshaler     = &Uint256{}
	_ yaml.Unmarshaler     = &Uint256{}
)

// Uint256 represents a nullable [bigutil.Uint256].
type Uint256 struct {
	Uint256 bigutil.Uint256
	Valid   bool
}

// NewUint256 returns a new [Uint256].
func NewUint256(x256 bigutil.Uint256, valid bool) Uint256 {
	return Uint256{
		Uint256: x256,
		Valid:   valid,
	}
}

// NullableString returns the string encoding of the underlying [bigutil.Uint256] as a [String], or an invalid [String] if n is invalid.
func (n Uint256) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Uint256.String(), true)
}

// Value implements [driver.Valuer].
// It encodes n by delegating to [bigutil.Uint256.Value] (or as nil if n is invalid).
func (n Uint256) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint256.Value()
}

// Scan implements [sql.Scanner].
// It decodes src into n by delegating to [bigutil.Uint256.Scan]; nil makes n invalid.
func (n *Uint256) Scan(src any) error {
	if src == nil {
		n.Uint256, n.Valid = bigutil.Uint256{}, false

		return nil
	}

	if err := n.Uint256.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n by delegating to [bigutil.Uint256.MarshalJSONTo] (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Uint256) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return n.Uint256.MarshalJSONTo(enc)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Uint256.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n by delegating to [bigutil.Uint256.MarshalYAML] (or as nil if n is invalid).
func (n Uint256) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint256.MarshalYAML()
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a value from dec into n by delegating to [bigutil.Uint256.UnmarshalJSONFrom]; the unquoted string null makes n invalid.
func (n *Uint256) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Uint256, n.Valid = bigutil.Uint256{}, false

		return nil

	default:
		var x256 bigutil.Uint256
		if err := x256.UnmarshalJSONFrom(dec); err != nil {
			return err
		}

		n.Uint256, n.Valid = x256, true

		return nil
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Uint256.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Uint256) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes value into n by delegating to [bigutil.Uint256.UnmarshalYAML]; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Uint256) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Uint256, n.Valid = bigutil.Uint256{}, false

		return nil
	}

	var x256 bigutil.Uint256
	if err := x256.UnmarshalYAML(value); err != nil {
		return err
	}

	n.Uint256, n.Valid = x256, true

	return nil
}
