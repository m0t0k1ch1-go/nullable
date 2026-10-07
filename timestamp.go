package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/m0t0k1ch1-go/timeutil/v5"
	"go.yaml.in/yaml/v3"
)

var (
	_ driver.Valuer        = Timestamp{}
	_ sql.Scanner          = &Timestamp{}
	_ json.MarshalerTo     = Timestamp{}
	_ json.Marshaler       = Timestamp{}
	_ yaml.Marshaler       = Timestamp{}
	_ json.UnmarshalerFrom = &Timestamp{}
	_ json.Unmarshaler     = &Timestamp{}
	_ yaml.Unmarshaler     = &Timestamp{}
)

// Timestamp represents a nullable [timeutil.Timestamp].
type Timestamp struct {
	Timestamp timeutil.Timestamp
	Valid     bool
}

// NewTimestamp returns a new [Timestamp].
func NewTimestamp(ts timeutil.Timestamp, valid bool) Timestamp {
	return Timestamp{
		Timestamp: ts,
		Valid:     valid,
	}
}

// NullableString returns the string encoding of the underlying [timeutil.Timestamp] as a [String], or an invalid [String] if n is invalid.
func (n Timestamp) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Timestamp.String(), true)
}

// Value implements [driver.Valuer].
// It encodes n by delegating to [timeutil.Timestamp.Value] (or as nil if n is invalid).
func (n Timestamp) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Timestamp.Value()
}

// Scan implements [sql.Scanner].
// It decodes src into n by delegating to [timeutil.Timestamp.Scan]; nil makes n invalid.
func (n *Timestamp) Scan(src any) error {
	if src == nil {
		n.Timestamp, n.Valid = timeutil.Timestamp{}, false

		return nil
	}

	if err := n.Timestamp.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n by delegating to [timeutil.Timestamp.MarshalJSONTo] (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Timestamp) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return n.Timestamp.MarshalJSONTo(enc)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Timestamp.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n by delegating to [timeutil.Timestamp.MarshalYAML] (or as nil if n is invalid).
func (n Timestamp) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Timestamp.MarshalYAML()
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a value from dec into n by delegating to [timeutil.Timestamp.UnmarshalJSONFrom]; the unquoted string null makes n invalid.
func (n *Timestamp) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Timestamp, n.Valid = timeutil.Timestamp{}, false

		return nil

	default:
		var ts timeutil.Timestamp
		if err := ts.UnmarshalJSONFrom(dec); err != nil {
			return err
		}

		n.Timestamp, n.Valid = ts, true

		return nil
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Timestamp.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Timestamp) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes value into n by delegating to [timeutil.Timestamp.UnmarshalYAML]; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Timestamp) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Timestamp, n.Valid = timeutil.Timestamp{}, false

		return nil
	}

	var ts timeutil.Timestamp
	if err := ts.UnmarshalYAML(value); err != nil {
		return err
	}

	n.Timestamp, n.Valid = ts, true

	return nil
}
