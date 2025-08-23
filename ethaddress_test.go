package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestEthAddress(t *testing.T) {
	var n nullable.EthAddress
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestEthAddress_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			want nullable.String
		}{
			{
				"null",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				nullable.NewString("", false),
			},
			{
				"vitalik.eth",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				nullable.NewString("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", true),
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

func TestEthAddress_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			want driver.Value
		}{
			{
				"null",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				nil,
			},
			{
				"vitalik.eth",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				ethhexutil.MustDecode("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
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

func TestEthAddressScan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"string",
				"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
				"",
			},
			{
				"[]byte: empty",
				[]byte{},
				"",
			},
			{
				"[]byte: invalid",
				[]byte{0x00},
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.EthAddress
		}{
			{
				"nil",
				nil,
				nullable.NewEthAddress(ethcommon.Address{}, false),
			},
			{
				"[]byte: vitalik.eth",
				ethhexutil.MustDecode("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.EthAddress, n.EthAddress)
			})
		}
	})
}

func TestEthAddress_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			want []byte
		}{
			{
				"null",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				[]byte(`null`),
			},
			{
				"vitalik.eth",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
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

func TestEthAddress_UnmarshalJSON(t *testing.T) {
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
				"string: empty",
				[]byte(`""`),
				"",
			},
			{
				"string: invalid",
				[]byte(`"0x00"`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.EthAddress
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewEthAddress(ethcommon.Address{}, false),
			},
			{
				"string: vitalik.eth",
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.EthAddress, n.EthAddress)
			})
		}
	})
}
