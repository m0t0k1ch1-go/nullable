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

func TestEthHash(t *testing.T) {
	var n nullable.EthHash
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestEthHash_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			want nullable.String
		}{
			{
				"null",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				nullable.NewString("", false),
			},
			{
				"Bitcoin genesis block",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				nullable.NewString("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f", true),
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

func TestEthHash_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			want driver.Value
		}{
			{
				"null",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				nil,
			},
			{
				"Bitcoin genesis block",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				ethhexutil.MustDecode("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
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

func TestEthHash_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"string",
				"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
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
				var n nullable.EthHash
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.EthHash
		}{
			{
				"nil",
				nil,
				nullable.NewEthHash(ethcommon.Hash{}, false),
			},
			{
				"[]byte: Bitcoin genesis block",
				ethhexutil.MustDecode("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthHash
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.EthHash, n.EthHash)
			})
		}
	})
}

func TestEthHash_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			want []byte
		}{
			{
				"null",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				[]byte(`null`),
			},
			{
				"Bitcoin genesis block",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
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

func TestEthHash_UnmarshalJSON(t *testing.T) {
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
				var n nullable.EthHash
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.EthHash
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewEthHash(ethcommon.Hash{}, false),
			},
			{
				"string: Bitcoin genesis block",
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthHash
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.EthHash, n.EthHash)
			})
		}
	})
}
