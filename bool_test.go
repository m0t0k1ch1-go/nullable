package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/coreutil"

	"github.com/m0t0k1ch1-go/nullable/v2"
	"github.com/m0t0k1ch1-go/nullable/v2/internal/testutil"
)

func TestNewBoolFromPtr(t *testing.T) {
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
				coreutil.Ptr(true),
				nullable.NewBool(true, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, nullable.NewBoolFromPtr(tc.in))
			})
		}
	})
}

func TestBoolPtr(t *testing.T) {
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
				coreutil.Ptr(true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, tc.in.Ptr())
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
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, b)
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
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
