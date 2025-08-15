package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestNewBoolFromBoolPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *bool
			out  nullable.Bool
		}{
			{
				"nil",
				nil,
				nullable.NewBool(false, false),
			},
			{
				"not nil",
				lo.ToPtr(true),
				nullable.NewBool(true, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, nullable.NewBoolFromBoolPtr(tc.in))
			})
		}
	})
}

func TestBoolBoolPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			out  *bool
		}{
			{
				"nil",
				nullable.NewBool(false, false),
				nil,
			},
			{
				"not nil",
				nullable.NewBool(true, true),
				lo.ToPtr(true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, tc.in.BoolPtr())
			})
		}
	})
}

func TestBoolMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			out  []byte
		}{
			{
				"null",
				nullable.NewBool(false, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewBool(true, true),
				[]byte("true"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.Nil(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestBoolUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Bool
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewBool(false, false),
			},
			{
				"not null",
				[]byte("true"),
				nullable.NewBool(true, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Bool
				require.Nil(t, json.Unmarshal(tc.in, &n))

				require.Equal(t, tc.out, n)
			})
		}
	})
}
