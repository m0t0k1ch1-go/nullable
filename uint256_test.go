package nullable_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"strings"
	"testing"

	"github.com/m0t0k1ch1-go/bigutil/v3"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

var (
	maxUint256 = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
)

func TestUint256(t *testing.T) {
	var n nullable.Uint256
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestUint256_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			want nullable.String
		}{
			{
				"null",
				nullable.NewUint256(bigutil.Uint256{}, false),
				nullable.NewString("", false),
			},
			{
				"zero",
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
				nullable.NewString("0x0", true),
			},
			{
				"one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				nullable.NewString("0x1", true),
			},
			{
				"max",
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
				nullable.NewString("0x"+strings.Repeat("f", 64), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := tc.in.NullableString()
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.String, n.String)
			})
		}
	})
}

func TestUint256_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			want driver.Value
		}{
			{
				"null",
				nullable.NewUint256(bigutil.Uint256{}, false),
				nil,
			},
			{
				"zero",
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
				[]byte{0x0},
			},
			{
				"one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				[]byte{0x1},
			},
			{
				"max",
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
				bytes.Repeat([]byte{0xff}, 32),
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

func TestUint256_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"int64",
				int64(0),
				"",
			},
			{
				"[]byte: empty",
				[]byte{},
				"",
			},
			{
				"[]byte: exceeds 256 bits",
				append([]byte{0x01}, bytes.Repeat([]byte{0x00}, 32)...),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.Uint256
		}{
			{
				"nil",
				nil,
				nullable.NewUint256(bigutil.Uint256{}, false),
			},
			{
				"[]byte: zero",
				[]byte{0x00},
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"[]byte: one",
				[]byte{0x01},
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"[]byte: max",
				bytes.Repeat([]byte{0xff}, 32),
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint256.String(), n.Uint256.String())
			})
		}
	})
}

func TestUint256_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			want []byte
		}{
			{
				"null",
				nullable.NewUint256(bigutil.Uint256{}, false),
				[]byte(`null`),
			},
			{
				"zero",
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
				[]byte(`"0x0"`),
			},
			{
				"one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				[]byte(`"0x1"`),
			},
			{
				"max",
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
				[]byte(`"0x` + strings.Repeat("f", 64) + `"`),
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

func TestUint256_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"empty",
				[]byte{},
				"",
			},
			{
				"number: negative",
				[]byte(`-1`),
				"",
			},
			{
				"number: exceeds 256 bits",
				[]byte(`115792089237316195423570985008687907853269984665640564039457584007913129639936`),
				"",
			},
			{
				"number: fractional",
				[]byte(`0.0`),
				"",
			},
			{
				"number: exponential",
				[]byte(`0e0`),
				"",
			},
			{
				"string: empty",
				[]byte(`""`),
				"",
			},
			{
				"string: invalid decimal",
				[]byte(`"invalid"`),
				"",
			},
			{
				"string: negative decimal",
				[]byte(`"-1"`),
				"",
			},
			{
				"string: missing hex digits after 0x prefix",
				[]byte(`"0x"`),
				"",
			},
			{
				"string: hex contains invalid escape sequences",
				[]byte(`"0x\x"`),
				"",
			},
			{
				"string: hex contains non-hex characters",
				[]byte(`"0xg"`),
				"",
			},
			{
				"string: hex exceeds 256 bits",
				[]byte(`"0x1` + strings.Repeat("0", 64) + `"`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Uint256
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewUint256(bigutil.Uint256{}, false),
			},
			{
				"number: zero",
				[]byte(`0`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"number: one",
				[]byte(`1`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"number: max",
				[]byte(`115792089237316195423570985008687907853269984665640564039457584007913129639935`),
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
			},
			{
				"string: decimal zero",
				[]byte(`"0"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"string: decimal one",
				[]byte(`"1"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"string: decimal max",
				[]byte(`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`),
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
			},
			{
				"string: hex zero",
				[]byte(`"0x0"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"string: hex one",
				[]byte(`"0x1"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"string: mixedcase hex max",
				[]byte(`"0x` + strings.Repeat("fF", 32) + `"`),
				nullable.NewUint256(bigutil.MustNewUint256(maxUint256), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint256.String(), n.Uint256.String())
			})
		}
	})
}
