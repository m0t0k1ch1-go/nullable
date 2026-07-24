package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/m0t0k1ch1-go/urlutil"
)

// HTTPURL represents a nullable urlutil.HTTPURL.
type HTTPURL struct {
	HTTPURL urlutil.HTTPURL
	Valid   bool
}

// NewHTTPURL returns a new HTTPURL.
func NewHTTPURL(hu urlutil.HTTPURL, valid bool) HTTPURL {
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
// It returns the driver.Value returned by urlutil.HTTPURL.Value, or nil if invalid.
func (n HTTPURL) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.HTTPURL.Value()
}

// Scan implements sql.Scanner.
// It accepts any value supported by urlutil.HTTPURL.Scan, or nil.
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

// MarshalJSON implements json.Marshaler.
// It returns the JSON encoding of urlutil.HTTPURL, or null if invalid.
func (n HTTPURL) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.HTTPURL)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts the JSON value supported by urlutil.HTTPURL, or null.
func (n *HTTPURL) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.HTTPURL, n.Valid = urlutil.HTTPURL{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.HTTPURL); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
