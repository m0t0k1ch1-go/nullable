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

func TestFloat64(t *testing.T) {
	var n nullable.Float64
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestNewFloat64FromFloat64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *float64
			want nullable.Float64
		}{
			{
				"nil",
				nil,
				nullable.NewFloat64(0, false),
			},
			{
				"zero",
				ptr(float64(0)),
				nullable.NewFloat64(0, true),
			},
			{
				"smallest non-zero",
				ptr(math.SmallestNonzeroFloat64),
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
			},
			{
				"min",
				ptr(-math.MaxFloat64),
				nullable.NewFloat64(-math.MaxFloat64, true),
			},
			{
				"max",
				ptr(math.MaxFloat64),
				nullable.NewFloat64(math.MaxFloat64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewFloat64FromFloat64Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Float64, n.Float64)
			})
		}
	})

	t.Run("success: captures value at call time", func(t *testing.T) {
		f := ptr(math.Phi)
		n := nullable.NewFloat64FromFloat64Ptr(f)

		*f = 0

		require.True(t, n.Valid)
		require.Equal(t, math.Phi, n.Float64)
	})
}

func TestFloat64_Float64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Float64
			want *float64
		}{
			{
				"null",
				nullable.NewFloat64(0, false),
				nil,
			},
			{
				"zero",
				nullable.NewFloat64(0, true),
				ptr(float64(0)),
			},
			{
				"smallest non-zero",
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
				ptr(math.SmallestNonzeroFloat64),
			},
			{
				"min",
				nullable.NewFloat64(-math.MaxFloat64, true),
				ptr(-math.MaxFloat64),
			},
			{
				"max",
				nullable.NewFloat64(math.MaxFloat64, true),
				ptr(math.MaxFloat64),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				f := tc.in.Float64Ptr()
				require.Equal(t, tc.want, f)
			})
		}
	})

	t.Run("success: pointer refers to a copy", func(t *testing.T) {
		n := nullable.NewFloat64(math.Phi, true)
		f := n.Float64Ptr()

		*f = 0

		require.True(t, n.Valid)
		require.Equal(t, math.Phi, n.Float64)
	})
}

func TestFloat64_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Float64
			want []byte
		}{
			{
				"null",
				nullable.NewFloat64(0, false),
				[]byte(`null`),
			},
			{
				"zero",
				nullable.NewFloat64(0, true),
				[]byte(`0`),
			},
			{
				"smallest non-zero",
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
				[]byte(`5e-324`),
			},
			{
				"min",
				nullable.NewFloat64(-math.MaxFloat64, true),
				[]byte(`-1.7976931348623157e+308`),
			},
			{
				"max",
				nullable.NewFloat64(math.MaxFloat64, true),
				[]byte(`1.7976931348623157e+308`),
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

func TestFloat64_UnmarshalJSON(t *testing.T) {
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
				"number: less than min",
				[]byte(`-1e+309`),
				"",
			},
			{
				"number: greater than max",
				[]byte(`1e+309`),
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
				var n nullable.Float64
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Float64
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewFloat64(0, false),
			},
			{
				"zero",
				[]byte(`0`),
				nullable.NewFloat64(0, true),
			},
			{
				"smallest non-zero",
				[]byte(`5e-324`),
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
			},
			{
				"min",
				[]byte(`-1.7976931348623157e+308`),
				nullable.NewFloat64(-math.MaxFloat64, true),
			},
			{
				"max",
				[]byte(`1.7976931348623157e+308`),
				nullable.NewFloat64(math.MaxFloat64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Float64
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Float64, n.Float64)
			})
		}
	})
}
