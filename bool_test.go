package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestBool(t *testing.T) {
	var n nullable.Bool
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestNewBoolFromBoolPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *bool
			want nullable.Bool
		}{
			{
				"nil",
				nil,
				nullable.NewBool(false, false),
			},
			{
				"true",
				ptr(true),
				nullable.NewBool(true, true),
			},
			{
				"false",
				ptr(false),
				nullable.NewBool(false, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewBoolFromBoolPtr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Bool, n.Bool)
			})
		}
	})

	t.Run("success: captures value at call time", func(t *testing.T) {
		b := ptr(true)
		n := nullable.NewBoolFromBoolPtr(b)

		*b = false

		require.True(t, n.Valid)
		require.True(t, n.Bool)
	})
}

func TestBool_BoolPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			want *bool
		}{
			{
				"null",
				nullable.NewBool(false, false),
				nil,
			},
			{
				"true",
				nullable.NewBool(true, true),
				ptr(true),
			},
			{
				"false",
				nullable.NewBool(false, true),
				ptr(false),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b := tc.in.BoolPtr()
				require.Equal(t, tc.want, b)
			})
		}
	})

	t.Run("success: pointer refers to a copy", func(t *testing.T) {
		n := nullable.NewBool(true, true)
		b := n.BoolPtr()

		*b = false

		require.True(t, n.Valid)
		require.True(t, n.Bool)
	})
}

func TestBool_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			want []byte
		}{
			{
				"null",
				nullable.NewBool(false, false),
				[]byte(`null`),
			},
			{
				"true",
				nullable.NewBool(true, true),
				[]byte(`true`),
			},
			{
				"false",
				nullable.NewBool(false, true),
				[]byte(`false`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := tc.in.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestBool_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"number",
				[]byte(`0`),
				"",
			},
			{
				"string",
				[]byte(`"true"`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Bool
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Bool
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewBool(false, false),
			},
			{
				"boolean: true",
				[]byte(`true`),
				nullable.NewBool(true, true),
			},
			{
				"boolean: false",
				[]byte(`false`),
				nullable.NewBool(false, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Bool
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Bool, n.Bool)
			})
		}
	})
}
