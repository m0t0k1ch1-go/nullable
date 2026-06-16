package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/m0t0k1ch1-go/sqlutil/v3"
)

// HTTPURL represents a nullable sqlutil.HTTPURL.
type HTTPURL struct {
	HTTPURL sqlutil.HTTPURL
	Valid   bool
}

// NewHTTPURL returns a new HTTPURL.
func NewHTTPURL(hu sqlutil.HTTPURL, valid bool) HTTPURL {
	return HTTPURL{
		HTTPURL: hu,
		Valid:   valid,
	}
}

// NullableString returns the value as a String.
func (n HTTPURL) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.HTTPURL.String(), true)
}

// Value implements driver.Valuer.
// It returns the driver.Value returned by sqlutil.HTTPURL.Value, or nil if invalid.
func (n HTTPURL) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.HTTPURL.Value()
}

// Scan implements sql.Scanner.
// It accepts any value supported by sqlutil.HTTPURL.Scan, or nil.
func (n *HTTPURL) Scan(src any) error {
	if src == nil {
		n.HTTPURL, n.Valid = sqlutil.HTTPURL{}, false

		return nil
	}

	if err := n.HTTPURL.Scan(src); err != nil {
		return err
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements json.Marshaler.
// It returns the JSON encoding of sqlutil.HTTPURL, or null if invalid.
func (n HTTPURL) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.HTTPURL)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts the JSON value supported by sqlutil.HTTPURL, or null.
func (n *HTTPURL) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.HTTPURL, n.Valid = sqlutil.HTTPURL{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.HTTPURL); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
