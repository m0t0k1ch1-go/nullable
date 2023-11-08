package nullable

import (
	"bytes"
	"database/sql/driver"

	"github.com/m0t0k1ch1-go/timeutil/v3"
)

// Timestamp is a nullable github.com/m0t0k1ch1-go/timeutil.Timestamp.
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

// NullableString returns the String.
func (ts Timestamp) NullableString() String {
	if !ts.Valid {
		return NewString("", false)
	}

	return NewString(ts.Timestamp.String(), true)
}

// Value implements the driver.Valuer interface.
func (ts Timestamp) Value() (driver.Value, error) {
	if !ts.Valid {
		return nil, nil
	}

	return ts.Timestamp.Value()
}

// Scan implements the sql.Scanner interface.
func (ts *Timestamp) Scan(src any) error {
	if src == nil {
		ts.Timestamp, ts.Valid = timeutil.Timestamp{}, false

		return nil
	}

	if err := ts.Timestamp.Scan(src); err != nil {
		return err
	}

	ts.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (ts Timestamp) MarshalJSON() ([]byte, error) {
	if !ts.Valid {
		return []byte("null"), nil
	}

	return ts.Timestamp.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (ts *Timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		ts.Timestamp, ts.Valid = timeutil.Timestamp{}, false

		return nil
	}

	if err := ts.Timestamp.UnmarshalJSON(b); err != nil {
		return err
	}

	ts.Valid = true

	return nil
}
