package nullable_test

import (
	"database/sql/driver"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestNewUint64FromUint64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *uint64
			want nullable.Uint64
		}{
			{
				"nil",
				nil,
				nullable.NewUint64(0, false),
			},
			{
				"zero",
				ptr(uint64(0)),
				nullable.NewUint64(0, true),
			},
			{
				"one",
				ptr(uint64(1)),
				nullable.NewUint64(1, true),
			},
			{
				"max",
				ptr(uint64(math.MaxUint64)),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewUint64FromUint64Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint64, n.Uint64)
			})
		}
	})
}

func TestUint64_Uint64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want *uint64
		}{
			{
				"null",
				nullable.NewUint64(0, false),
				nil,
			},
			{
				"zero",
				nullable.NewUint64(0, true),
				ptr(uint64(0)),
			},
			{
				"one",
				nullable.NewUint64(1, true),
				ptr(uint64(1)),
			},
			{
				"max",
				nullable.NewUint64(math.MaxUint64, true),
				ptr(uint64(math.MaxUint64)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := tc.in
				p := n.Uint64Ptr()
				require.Equal(t, tc.want, p)

				if p != nil {
					*p = math.MaxUint8

					require.Equal(t, tc.in.Valid, n.Valid)
					require.Equal(t, tc.in.Uint64, n.Uint64)
				}
			})
		}
	})
}

func TestUint64_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want driver.Value
		}{
			{
				"null",
				nullable.NewUint64(0, false),
				nil,
			},
			{
				"zero",
				nullable.NewUint64(0, true),
				uint64(0),
			},
			{
				"one",
				nullable.NewUint64(1, true),
				uint64(1),
			},
			{
				"max",
				nullable.NewUint64(math.MaxUint64, true),
				uint64(math.MaxUint64),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestUint64_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"float64",
				float64(0),
				"unsupported source type: float64",
			},
			{
				"int64: negative",
				int64(-1),
				"invalid source: negative int64",
			},
			{
				"[]byte: empty",
				[]byte{},
				"invalid source: empty []byte",
			},
			{
				"[]byte: invalid",
				[]byte("invalid"),
				"invalid source",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.Uint64
		}{
			{
				"nil",
				nil,
				nullable.NewUint64(0, false),
			},
			{
				"int64: zero",
				int64(0),
				nullable.NewUint64(0, true),
			},
			{
				"int64: one",
				int64(1),
				nullable.NewUint64(1, true),
			},
			{
				"int64: max",
				int64(math.MaxInt64),
				nullable.NewUint64(math.MaxInt64, true),
			},
			{
				"uint64: zero",
				uint64(0),
				nullable.NewUint64(0, true),
			},
			{
				"uint64: one",
				uint64(1),
				nullable.NewUint64(1, true),
			},
			{
				"uint64: max",
				uint64(math.MaxUint64),
				nullable.NewUint64(math.MaxUint64, true),
			},
			{
				"[]byte: zero",
				[]byte("0"),
				nullable.NewUint64(0, true),
			},
			{
				"[]byte: one",
				[]byte("1"),
				nullable.NewUint64(1, true),
			},
			{
				"[]byte: max",
				[]byte("18446744073709551615"),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint64, n.Uint64)
			})
		}
	})
}

func TestUint64_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want []byte
		}{
			{
				"null",
				nullable.NewUint64(0, false),
				[]byte(`null`),
			},
			{
				"zero",
				nullable.NewUint64(0, true),
				[]byte(`0`),
			},
			{
				"one",
				nullable.NewUint64(1, true),
				[]byte(`1`),
			},
			{
				"max",
				nullable.NewUint64(math.MaxUint64, true),
				[]byte(`18446744073709551615`),
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

func TestUint64_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"boolean",
				[]byte(`true`),
				"",
			},
			{
				"string",
				[]byte(`"0"`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Uint64
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewUint64(0, false),
			},
			{
				"number: zero",
				[]byte(`0`),
				nullable.NewUint64(0, true),
			},
			{
				"number: one",
				[]byte(`1`),
				nullable.NewUint64(1, true),
			},
			{
				"number: max",
				[]byte(`18446744073709551615`),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint64, n.Uint64)
			})
		}
	})
}
