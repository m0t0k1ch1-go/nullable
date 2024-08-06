package nullable_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/m0t0k1ch1-go/nullable/v2"
	"github.com/m0t0k1ch1-go/nullable/v2/internal/testutil"
)

func TestEthAddressNullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			out  nullable.String
		}{
			{
				"null",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				nullable.NewString("", false),
			},
			{
				"not null",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				nullable.NewString("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, tc.in.NullableString())
			})
		}
	})
}

func TestEthAddressValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			out  driver.Value
		}{
			{
				"null",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				nil,
			},
			{
				"not null",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				ethhexutil.MustDecode("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
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

func TestEthAddressScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  nullable.EthAddress
		}{
			{
				"null",
				nil,
				nullable.NewEthAddress(ethcommon.Address{}, false),
			},
			{
				"not null",
				ethhexutil.MustDecode("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				if err := n.Scan(tc.in); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}

func TestEthAddressMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			out  []byte
		}{
			{
				"null",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
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

func TestEthAddressUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.EthAddress
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewEthAddress(ethcommon.Address{}, false),
			},
			{
				"not null",
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
