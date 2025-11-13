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

func TestInt32(t *testing.T) {
	var n nullable.Int32
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestNewInt32FromInt32Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *int32
			want nullable.Int32
		}{
			{
				"nil",
				nil,
				nullable.NewInt32(0, false),
			},
			{
				"zero",
				ptr(int32(0)),
				nullable.NewInt32(0, true),
			},
			{
				"min",
				ptr(int32(math.MinInt32)),
				nullable.NewInt32(math.MinInt32, true),
			},
			{
				"max",
				ptr(int32(math.MaxInt32)),
				nullable.NewInt32(math.MaxInt32, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewInt32FromInt32Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int32, n.Int32)
			})
		}
	})

	t.Run("success: captures value at call time", func(t *testing.T) {
		i := ptr(int32(1))
		n := nullable.NewInt32FromInt32Ptr(i)

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int32(1), n.Int32)
	})
}

func TestInt32_Int32Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			want *int32
		}{
			{
				"null",
				nullable.NewInt32(0, false),
				nil,
			},
			{
				"zero",
				nullable.NewInt32(0, true),
				ptr(int32(0)),
			},
			{
				"min",
				nullable.NewInt32(math.MinInt32, true),
				ptr(int32(math.MinInt32)),
			},
			{
				"max",
				nullable.NewInt32(math.MaxInt32, true),
				ptr(int32(math.MaxInt32)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				i := tc.in.Int32Ptr()
				require.Equal(t, tc.want, i)
			})
		}
	})

	t.Run("success: pointer refers to a copy", func(t *testing.T) {
		n := nullable.NewInt32(1, true)
		i := n.Int32Ptr()

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int32(1), n.Int32)
	})
}

func TestInt32_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			want []byte
		}{
			{
				"null",
				nullable.NewInt32(0, false),
				[]byte(`null`),
			},
			{
				"zero",
				nullable.NewInt32(0, true),
				[]byte(`0`),
			},
			{
				"min",
				nullable.NewInt32(math.MinInt32, true),
				[]byte(`-2147483648`),
			},
			{
				"max",
				nullable.NewInt32(math.MaxInt32, true),
				[]byte(`2147483647`),
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

func TestInt32_UnmarshalJSON(t *testing.T) {
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
				[]byte(`-2147483649`),
				"",
			},
			{
				"number: max + 1",
				[]byte(`2147483648`),
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
				var n nullable.Int32
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Int32
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewInt32(0, false),
			},
			{
				"number: zero",
				[]byte(`0`),
				nullable.NewInt32(0, true),
			},
			{
				"number: min",
				[]byte(`-2147483648`),
				nullable.NewInt32(math.MinInt32, true),
			},
			{
				"number: max",
				[]byte(`2147483647`),
				nullable.NewInt32(math.MaxInt32, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int32
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int32, n.Int32)
			})
		}
	})
}
