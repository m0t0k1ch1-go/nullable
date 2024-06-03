package nullable_test

import (
	"encoding/json"
	"testing"

	"database/sql/driver"

	"github.com/m0t0k1ch1-go/nullable"
	"github.com/m0t0k1ch1-go/nullable/internal/testutil"
)

func TestUint64Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			out  driver.Value
		}{
			{
				"null",
				nullable.NewUint64(0, false),
				nil,
			},
			{
				"not null",
				nullable.NewUint64(1231006505, true),
				uint64(1231006505),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, v)
			})
		}
	})
}

func TestUint64Scan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  nullable.Uint64
		}{
			{
				"null",
				nil,
				nullable.NewUint64(0, false),
			},
			{
				"positive int64",
				int64(1231006505),
				nullable.NewUint64(1231006505, true),
			},
			{
				"uint64",
				uint64(1231006505),
				nullable.NewUint64(1231006505, true),
			},
			{
				"[]byte",
				[]byte("1231006505"),
				nullable.NewUint64(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				if err := n.Scan(tc.in); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}

func TestUint64MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			out  []byte
		}{
			{
				"null",
				nullable.NewUint64(0, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewUint64(1231006505, true),
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

func TestUint64UnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Uint64
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewUint64(0, false),
			},
			{
				"not null",
				[]byte("1231006505"),
				nullable.NewUint64(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
