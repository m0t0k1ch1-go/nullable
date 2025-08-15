package nullable_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	ethmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/m0t0k1ch1-go/bigutil/v3"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestUint256NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			out  nullable.String
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
				"max",
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
				nullable.NewString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := tc.in.NullableString()

				require.Equal(t, tc.out, n)
			})
		}
	})
}

func TestUint256Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			out  driver.Value
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
				"max",
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				require.NoError(t, err)

				require.Equal(t, tc.out, v)
			})
		}
	})
}

func TestUint256Scan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  nullable.Uint256
		}{
			{
				"nil",
				nil,
				nullable.NewUint256(bigutil.Uint256{}, false),
			},
			{
				"zero",
				[]byte{0x0},
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"max",
				[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				{
					err := n.Scan(tc.in)
					require.NoError(t, err)
				}

				require.Equal(t, tc.out.Valid, n.Valid)
				require.Zero(t, n.Uint256.BigInt().Cmp(tc.out.Uint256.BigInt()))
			})
		}
	})
}

func TestUint256MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			out  []byte
		}{
			{
				"null",
				nullable.NewUint256(bigutil.Uint256{}, false),
				[]byte("null"),
			},
			{
				"zero",
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
				[]byte(`"0x0"`),
			},
			{
				"max",
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.NoError(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestUint256UnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Uint256
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewUint256(bigutil.Uint256{}, false),
			},
			{
				"zero (hexadecimal string)",
				[]byte(`"0x0"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"max (hexadecimal string)",
				[]byte(`"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`),
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
			},
			{
				"zero (decimal string)",
				[]byte(`"0"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"max (decimal string)",
				[]byte(`"115792089237316195423570985008687907853269984665640564039457584007913129639935"`),
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
			},
			{
				"zero (number)",
				[]byte("0"),
				nullable.NewUint256(bigutil.NewUint256FromUint64(0), true),
			},
			{
				"max (number)",
				[]byte("115792089237316195423570985008687907853269984665640564039457584007913129639935"),
				nullable.NewUint256(bigutil.MustNewUint256(ethmath.MaxBig256), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				{
					err := json.Unmarshal(tc.in, &n)
					require.NoError(t, err)
				}

				require.Equal(t, tc.out.Valid, n.Valid)
				require.Zero(t, n.Uint256.BigInt().Cmp(tc.out.Uint256.BigInt()))
			})
		}
	})
}
