package nullable_test

import (
	"encoding/json"
	"testing"

	"database/sql/driver"

	"github.com/m0t0k1ch1-go/coreutil"
	"github.com/m0t0k1ch1-go/nullable/v2"
	"github.com/stretchr/testify/require"
)

func TestNewUint64FromUint64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *uint64
			out  nullable.Uint64
		}{
			{
				"nil",
				nil,
				nullable.NewUint64(0, false),
			},
			{
				"not nil",
				coreutil.Ptr(uint64(1231006505)),
				nullable.NewUint64(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, nullable.NewUint64FromUint64Ptr(tc.in))
			})
		}
	})
}

func TestUint64Uint64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			out  *uint64
		}{
			{
				"nil",
				nullable.NewUint64(0, false),
				nil,
			},
			{
				"not nil",
				nullable.NewUint64(1231006505, true),
				coreutil.Ptr(uint64(1231006505)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, tc.in.Uint64Ptr())
			})
		}
	})
}

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
				require.Nil(t, err)

				require.Equal(t, tc.out, v)
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
				require.Nil(t, n.Scan(tc.in))

				require.Equal(t, tc.out, n)
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
				require.Nil(t, err)

				require.Equal(t, tc.out, b)
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
				require.Nil(t, json.Unmarshal(tc.in, &n))

				require.Equal(t, tc.out, n)
			})
		}
	})
}
