package nullable_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

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
				n := tc.in
				p := n.BoolPtr()
				require.Equal(t, tc.want, p)

				if p != nil {
					*p = !tc.in.Bool

					require.Equal(t, tc.in.Valid, n.Valid)
					require.Equal(t, tc.in.Bool, n.Bool)
				}
			})
		}
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
