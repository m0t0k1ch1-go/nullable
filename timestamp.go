package nullable

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/m0t0k1ch1-go/timeutil/v4"
	"github.com/samber/oops"
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
func (n Timestamp) NullableString() String {
	if !n.Valid {
		return NewString("", false)
	}

	return NewString(n.Timestamp.String(), true)
}

// Value implements the driver.Valuer interface.
func (n Timestamp) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	v, err := n.Timestamp.Value()
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return v, nil
}

// Scan implements the sql.Scanner interface.
func (n *Timestamp) Scan(src any) error {
	if src == nil {
		n.Timestamp, n.Valid = timeutil.Timestamp{}, false

		return nil
	}

	if err := n.Timestamp.Scan(src); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}

// MarshalJSON implements the json.Marshaler interface.
func (n Timestamp) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}

	b, err := json.Marshal(n.Timestamp)
	if err != nil {
		return nil, oops.Wrap(err)
	}

	return b, nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (n *Timestamp) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Timestamp, n.Valid = timeutil.Timestamp{}, false

		return nil
	}

	if err := json.Unmarshal(b, &n.Timestamp); err != nil {
		return oops.Wrap(err)
	}

	n.Valid = true

	return nil
}
