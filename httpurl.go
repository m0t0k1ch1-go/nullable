package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/m0t0k1ch1-go/urlutil"
	"go.yaml.in/yaml/v3"
)

var (
	_ driver.Valuer        = HTTPURL{}
	_ sql.Scanner          = &HTTPURL{}
	_ json.MarshalerTo     = HTTPURL{}
	_ json.Marshaler       = HTTPURL{}
	_ yaml.Marshaler       = HTTPURL{}
	_ json.UnmarshalerFrom = &HTTPURL{}
	_ json.Unmarshaler     = &HTTPURL{}
	_ yaml.Unmarshaler     = &HTTPURL{}
)

// HTTPURL represents a nullable [urlutil.HTTPURL].
type HTTPURL struct {
	HTTPURL urlutil.HTTPURL
	Valid   bool
}

// NewHTTPURL returns a new [HTTPURL].
func NewHTTPURL(hu urlutil.HTTPURL, valid bool) HTTPURL {
	return HTTPURL{
		HTTPURL: hu,
		Valid:   valid,
	}
}

// NullableString returns the string encoding of the underlying [urlutil.HTTPURL] as a [String], or an invalid [String] if n is invalid.
func (n HTTPURL) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.HTTPURL.String(), true)
}

// Value implements [driver.Valuer].
// It encodes n by delegating to [urlutil.HTTPURL.Value] (or as nil if n is invalid).
func (n HTTPURL) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.HTTPURL.Value()
}

// Scan implements [sql.Scanner].
// It decodes src into n by delegating to [urlutil.HTTPURL.Scan]; nil makes n invalid.
func (n *HTTPURL) Scan(src any) error {
	if src == nil {
		n.HTTPURL, n.Valid = urlutil.HTTPURL{}, false

		return nil
	}

	if err := n.HTTPURL.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n by delegating to [urlutil.HTTPURL.MarshalJSONTo] (or as the unquoted string null if n is invalid) and writes it to enc.
func (n HTTPURL) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return n.HTTPURL.MarshalJSONTo(enc)
}

// MarshalJSON implements [json.Marshaler].
// It is like [HTTPURL.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n HTTPURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n by delegating to [urlutil.HTTPURL.MarshalYAML] (or as nil if n is invalid).
func (n HTTPURL) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.HTTPURL.MarshalYAML()
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes a value from dec into n by delegating to [urlutil.HTTPURL.UnmarshalJSONFrom]; the unquoted string null makes n invalid.
func (n *HTTPURL) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.HTTPURL, n.Valid = urlutil.HTTPURL{}, false

		return nil

	default:
		var hu urlutil.HTTPURL
		if err := hu.UnmarshalJSONFrom(dec); err != nil {
			return err
		}

		n.HTTPURL, n.Valid = hu, true

		return nil
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [HTTPURL.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *HTTPURL) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes value into n by delegating to [urlutil.HTTPURL.UnmarshalYAML]; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *HTTPURL) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.HTTPURL, n.Valid = urlutil.HTTPURL{}, false

		return nil
	}

	var hu urlutil.HTTPURL
	if err := hu.UnmarshalYAML(value); err != nil {
		return err
	}

	n.HTTPURL, n.Valid = hu, true

	return nil
}
