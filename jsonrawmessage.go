package nullable

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/samber/oops"
)

// JSONRawMessage is a nullable json.RawMessage.
type JSONRawMessage struct {
	JSONRawMessage json.RawMessage
	Valid          bool
}

// NewJSONRawMessage returns a new JSONRawMessage.
func NewJSONRawMessage(msg json.RawMessage, valid bool) JSONRawMessage {
	return JSONRawMessage{
		JSONRawMessage: msg,
		Valid:          valid,
	}
}

// Value implements the driver.Valuer interface.
func (n JSONRawMessage) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}

	return []byte(n.JSONRawMessage), nil
}

// Scan implements the sql.Scanner interface.
func (n *JSONRawMessage) Scan(src any) error {
	if src == nil {
		n.JSONRawMessage, n.Valid = nil, false

		return nil
	}

	b, ok := src.([]byte)
	if !ok {
		return oops.Errorf("unexpected src type: %T", src)
	}

	n.JSONRawMessage, n.Valid = b, true

	return nil
}
