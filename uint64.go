package nullable

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"strconv"

	"go.yaml.in/yaml/v3"
)

var (
	_ driver.Valuer        = Uint64{}
	_ sql.Scanner          = &Uint64{}
	_ json.MarshalerTo     = Uint64{}
	_ json.Marshaler       = Uint64{}
	_ yaml.Marshaler       = Uint64{}
	_ json.UnmarshalerFrom = &Uint64{}
	_ json.Unmarshaler     = &Uint64{}
	_ yaml.Unmarshaler     = &Uint64{}
)

// Uint64 represents a nullable uint64.
type Uint64 struct {
	Uint64 uint64
	Valid  bool
}

// NewUint64 returns a new [Uint64].
func NewUint64(i uint64, valid bool) Uint64 {
	return Uint64{
		Uint64: i,
		Valid:  valid,
	}
}

// NewUint64FromUint64Ptr returns a new [Uint64] from a uint64 pointer.
// It returns an invalid [Uint64] if i is nil.
func NewUint64FromUint64Ptr(i *uint64) Uint64 {
	if i == nil {
		return NewUint64(0, false)
	}

	return NewUint64(*i, true)
}

// Uint64Ptr returns a pointer to a copy of the underlying uint64, or nil if n is invalid.
func (n Uint64) Uint64Ptr() *uint64 {
	if !n.Valid {
		return nil
	}

	return &n.Uint64
}

// Value implements [driver.Valuer].
// It encodes n as a uint64 (or nil if n is invalid).
func (n Uint64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint64, nil
}

// Scan implements [sql.Scanner].
// It decodes one of the following values into n; nil makes n invalid:
//   - nil
//   - int64 (non-negative)
//   - uint64
//   - decimal string bytes
func (n *Uint64) Scan(src any) error {
	if src == nil {
		n.Uint64, n.Valid = 0, false

		return nil
	}

	switch src := src.(type) {
	case int64:
		if src < 0 {
			return errors.New("invalid int64 source: negative")
		}

		n.Uint64, n.Valid = uint64(src), true

		return nil

	case uint64:
		n.Uint64, n.Valid = src, true

		return nil

	case []byte:
		if len(src) == 0 {
			return errors.New("invalid bytes source: empty")
		}

		i, err := strconv.ParseUint(string(src), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid bytes source: %w", err)
		}

		n.Uint64, n.Valid = i, true

		return nil

	default:
		return fmt.Errorf("unsupported source type: %T", src)
	}
}

// MarshalJSONTo implements [json.MarshalerTo].
// It encodes n as an unquoted decimal string (or as the unquoted string null if n is invalid) and writes it to enc.
func (n Uint64) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !n.Valid {
		return enc.WriteToken(jsontext.Null)
	}

	return json.MarshalEncode(enc, n.Uint64)
}

// MarshalJSON implements [json.Marshaler].
// It is like [Uint64.MarshalJSONTo] but returns the encoded bytes instead of writing them to a [jsontext.Encoder].
func (n Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(n)
}

// MarshalYAML implements [yaml.Marshaler].
// It encodes n as a uint64 (or nil if n is invalid).
func (n Uint64) MarshalYAML() (any, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Uint64, nil
}

// UnmarshalJSONFrom implements [json.UnmarshalerFrom].
// It decodes an unquoted decimal string or the unquoted string null from dec into n; the latter makes n invalid.
func (n *Uint64) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch k := dec.PeekKind(); k {
	case jsontext.KindNumber:
		var i uint64
		if err := json.UnmarshalDecode(dec, &i); err != nil {
			return fmt.Errorf("invalid number: %w", err)
		}

		n.Uint64, n.Valid = i, true

		return nil

	case jsontext.KindNull:
		if _, err := dec.ReadToken(); err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}

		n.Uint64, n.Valid = 0, false

		return nil

	default:
		return fmt.Errorf("unsupported json token kind: %v", k)
	}
}

// UnmarshalJSON implements [json.Unmarshaler].
// It is like [Uint64.UnmarshalJSONFrom] but decodes b instead of reading from a [jsontext.Decoder].
func (n *Uint64) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, n)
}

// UnmarshalYAML implements [yaml.Unmarshaler].
// It decodes a scalar from value into n as a uint64; a scalar tagged !!null makes n invalid.
// Note that [go.yaml.in/yaml/v3] never calls this method for null nodes and leaves n unchanged instead.
func (n *Uint64) UnmarshalYAML(value *yaml.Node) error {
	if value.ShortTag() == "!!null" {
		n.Uint64, n.Valid = 0, false

		return nil
	}

	var i uint64
	if err := value.Decode(&i); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	n.Uint64, n.Valid = i, true

	return nil
}
