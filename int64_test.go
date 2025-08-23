package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestInt64(t *testing.T) {
	var n nullable.Int64
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

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
				"zero",
				ptr(int64(0)),
				nullable.NewInt64(0, true),
			},
			{
				"min",
				ptr(int64(math.MinInt64)),
				nullable.NewInt64(math.MinInt64, true),
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

	t.Run("success: captures value at call time", func(t *testing.T) {
		i := ptr(int64(1))
		n := nullable.NewInt64FromInt64Ptr(i)

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int64(1), n.Int64)
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
				"zero",
				nullable.NewInt64(0, true),
				ptr(int64(0)),
			},
			{
				"min",
				nullable.NewInt64(math.MinInt64, true),
				ptr(int64(math.MinInt64)),
			},
			{
				"max",
				nullable.NewInt64(math.MaxInt64, true),
				ptr(int64(math.MaxInt64)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				i := tc.in.Int64Ptr()
				require.Equal(t, tc.want, i)
			})
		}
	})

	t.Run("success: pointer refers to a copy", func(t *testing.T) {
		n := nullable.NewInt64(1, true)
		i := n.Int64Ptr()

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int64(1), n.Int64)
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
				"zero",
				nullable.NewInt64(0, true),
				[]byte(`0`),
			},
			{
				"min",
				nullable.NewInt64(math.MinInt64, true),
				[]byte(`-9223372036854775808`),
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

func TestInt64_UnmarshalJSON(t *testing.T) {
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
				"number: min - 1",
				[]byte(`-9223372036854775809`),
				"",
			},
			{
				"number: max + 1",
				[]byte(`9223372036854775808`),
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
				"number: zero",
				[]byte(`0`),
				nullable.NewInt64(0, true),
			},
			{
				"number: min",
				[]byte(`-9223372036854775808`),
				nullable.NewInt64(math.MinInt64, true),
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
