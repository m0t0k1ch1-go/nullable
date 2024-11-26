package nullable_test

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestJSONRawMessageValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.JSONRawMessage
			out  driver.Value
		}{
			{
				"null",
				nullable.NewJSONRawMessage(nil, false),
				nil,
			},
			{
				"not null",
				nullable.NewJSONRawMessage([]byte(`{"k":"v"}`), true),
				[]byte(`{"k":"v"}`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.Nil(t, err)

				require.Equal(t, tc.out, v)
			})
		}
	})
}

func TestJSONRawMessageScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  nullable.JSONRawMessage
		}{
			{
				"null",
				nil,
				nullable.NewJSONRawMessage(nil, false),
			},
			{
				"not null",
				[]byte(`{"k":"v"}`),
				nullable.NewJSONRawMessage([]byte(`{"k":"v"}`), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.JSONRawMessage
				require.Nil(t, n.Scan(tc.in))

				require.Equal(t, tc.out, n)
			})
		}
	})
}
