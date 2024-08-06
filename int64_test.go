package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/coreutil"

	"github.com/m0t0k1ch1-go/nullable/v2"
	"github.com/m0t0k1ch1-go/nullable/v2/internal/testutil"
)

func TestNewInt64FromInt64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *int64
			out  nullable.Int64
		}{
			{
				"nil",
				nil,
				nullable.NewInt64(0, false),
			},
			{
				"not nil",
				coreutil.Ptr(int64(1231006505)),
				nullable.NewInt64(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, nullable.NewInt64FromInt64Ptr(tc.in))
			})
		}
	})
}

func TestInt64Int64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			out  *int64
		}{
			{
				"nil",
				nullable.NewInt64(0, false),
				nil,
			},
			{
				"not nil",
				nullable.NewInt64(1231006505, true),
				coreutil.Ptr(int64(1231006505)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, tc.in.Int64Ptr())
			})
		}
	})
}

func TestInt64MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			out  []byte
		}{
			{
				"null",
				nullable.NewInt64(0, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewInt64(1231006505, true),
				[]byte("1231006505"),
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

func TestInt64UnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Int64
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewInt64(0, false),
			},
			{
				"not null",
				[]byte("1231006505"),
				nullable.NewInt64(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int64
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
