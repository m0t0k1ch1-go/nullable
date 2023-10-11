package nullable

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// String is a nullable string.
type String struct {
	sql.NullString
}

// NewString returns a new String.
func NewString(s string, valid bool) String {
	return String{
		sql.NullString{
			String: s,
			Valid:  valid,
		},
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (s String) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(s.String)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (s *String) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		s.String, s.Valid = "", false

		return nil
	}

	if err := json.Unmarshal(b, &s.String); err != nil {
		return err
	}

	s.Valid = true

	return nil
}
