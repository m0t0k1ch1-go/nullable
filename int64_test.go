package nullable_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestNewInt64FromInt64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *int64
			want nullable.Int64
		}{
			{
				"nil",
				nil,
				nullable.NewInt64(0, false),
			},
			{
				"min",
				ptr(int64(math.MinInt64)),
				nullable.NewInt64(math.MinInt64, true),
			},
			{
				"zero",
				ptr(int64(0)),
				nullable.NewInt64(0, true),
			},
			{
				"max",
				ptr(int64(math.MaxInt64)),
				nullable.NewInt64(math.MaxInt64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewInt64FromInt64Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int64, n.Int64)
			})
		}
	})
}

func TestInt64_Int64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			want *int64
		}{
			{
				"null",
				nullable.NewInt64(0, false),
				nil,
			},
			{
				"min",
				nullable.NewInt64(math.MinInt64, true),
				ptr(int64(math.MinInt64)),
			},
			{
				"zero",
				nullable.NewInt64(0, true),
				ptr(int64(0)),
			},
			{
				"max",
				nullable.NewInt64(math.MaxInt64, true),
				ptr(int64(math.MaxInt64)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := tc.in
				p := n.Int64Ptr()
				require.Equal(t, tc.want, p)

				if p != nil {
					*p = 1

					require.Equal(t, tc.in.Valid, n.Valid)
					require.Equal(t, tc.in.Int64, n.Int64)
				}
			})
		}
	})
}

func TestInt64_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			want []byte
		}{
			{
				"null",
				nullable.NewInt64(0, false),
				[]byte(`null`),
			},
			{
				"min",
				nullable.NewInt64(math.MinInt64, true),
				[]byte(`-9223372036854775808`),
			},
			{
				"zero",
				nullable.NewInt64(0, true),
				[]byte(`0`),
			},
			{
				"max",
				nullable.NewInt64(math.MaxInt64, true),
				[]byte(`9223372036854775807`),
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

func TestInt64UnmarshalJSON(t *testing.T) {
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
				var n nullable.Int64
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Int64
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewInt64(0, false),
			},
			{
				"number: min",
				[]byte(`-9223372036854775808`),
				nullable.NewInt64(math.MinInt64, true),
			},
			{
				"number: zero",
				[]byte(`0`),
				nullable.NewInt64(0, true),
			},
			{
				"number: max",
				[]byte(`9223372036854775807`),
				nullable.NewInt64(math.MaxInt64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int64
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int64, n.Int64)
			})
		}
	})
}
