package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/m0t0k1ch1-go/timeutil/v5"
)

// Timestamp represents a nullable timeutil.Timestamp.
type Timestamp struct {
	Timestamp timeutil.Timestamp
	Valid     bool
}

// NewTimestamp returns a new Timestamp.
func NewTimestamp(ts timeutil.Timestamp, valid bool) Timestamp {
	return Timestamp{
		Timestamp: ts,
		Valid:     valid,
	}
}

// NullableString returns the value as a String.
func (n Timestamp) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Timestamp.String(), true)
}

// Value implements driver.Valuer.
// It returns the driver.Value representation of timeutil.Timestamp, or nil if invalid.
func (n Timestamp) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return n.Timestamp.Value()
}

// Scan implements sql.Scanner.
// It accepts a value supported by timeutil.Timestamp.Scan, or nil.
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

// MarshalJSON implements json.Marshaler.
// It returns the value as the JSON encoding of timeutil.Timestamp, or null if invalid.
func (n Timestamp) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.Timestamp)
}

// UnmarshalJSON implements json.Unmarshaler.
// It accepts any JSON value supported by timeutil.Timestamp, or null.
func (n *Timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Timestamp, n.Valid = timeutil.Timestamp{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Timestamp); err != nil {
		return err
	}

	n.Valid = true

	return nil
}
