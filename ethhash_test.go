package nullable_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestEthHashNullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			out  nullable.String
		}{
			{
				"null",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				nullable.NewString("", false),
			},
			{
				"not null",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				nullable.NewString("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f", true),
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

func TestEthHashValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			out  driver.Value
		}{
			{
				"null",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				nil,
			},
			{
				"not null",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				ethhexutil.MustDecode("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
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

func TestEthHashScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  nullable.EthHash
		}{
			{
				"nil",
				nil,
				nullable.NewEthHash(ethcommon.Hash{}, false),
			},
			{
				"not nil",
				ethhexutil.MustDecode("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthHash
				{
					err := n.Scan(tc.in)
					require.NoError(t, err)
				}

				require.Equal(t, tc.out, n)
			})
		}
	})
}

func TestEthHashMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			out  []byte
		}{
			{
				"null",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
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

func TestEthHashUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.EthHash
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewEthHash(ethcommon.Hash{}, false),
			},
			{
				"not null",
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthHash
				{
					err := json.Unmarshal(tc.in, &n)
					require.NoError(t, err)
				}

				require.Equal(t, tc.out, n)
			})
		}
	})
}
