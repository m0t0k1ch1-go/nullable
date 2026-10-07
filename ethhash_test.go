package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestEthHash(t *testing.T) {
	var n nullable.EthHash
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestEthHash_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			want nullable.String
		}{
			{
				"invalid",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				nullable.NewString("", false),
			},
			{
				"valid: zero",
				nullable.NewEthHash(ethcommon.Hash{}, true),
				nullable.NewString("0x0000000000000000000000000000000000000000000000000000000000000000", true),
			},
			{
				"valid: bitcoin genesis block",
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
				"invalid",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewEthHash(ethcommon.Hash{}, true),
				ethhexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000"),
			},
			{
				"valid: bitcoin genesis block",
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
				"bytes: zero",
				ethhexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000"),
				nullable.NewEthHash(ethcommon.Hash{}, true),
			},
			{
				"bytes: bitcoin genesis block",
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

func TestEthHash_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.EthHash) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.EthHash) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.EthHash) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			want []byte
		}{
			{
				"invalid",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				[]byte(`null`),
			},
			{
				"valid: zero",
				nullable.NewEthHash(ethcommon.Hash{}, true),
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000000"`),
			},
			{
				"valid: bitcoin genesis block",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, enc := range encs {
					t.Run(enc.name, func(t *testing.T) {
						b, err := enc.marshal(tc.in)
						require.NoError(t, err)
						require.Equal(t, tc.want, b)
					})
				}
			})
		}
	})
}

func TestEthHash_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthHash
			want []byte
		}{
			{
				"invalid",
				nullable.NewEthHash(ethcommon.Hash{}, false),
				[]byte("null\n"),
			},
			{
				"valid: zero",
				nullable.NewEthHash(ethcommon.Hash{}, true),
				[]byte("\"0x0000000000000000000000000000000000000000000000000000000000000000\"\n"),
			},
			{
				"valid: bitcoin genesis block",
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
				[]byte("\"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f\"\n"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := yaml.Marshal(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestEthHash_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.EthHash) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.EthHash) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.EthHash) error {
				return n.UnmarshalJSON(b)
			},
		},
	}

	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"nil",
				nil,
				"",
			},
			{
				"empty",
				[]byte{},
				"",
			},
			{
				"unquoted string bytes: boolean",
				[]byte(`true`),
				"unsupported json token kind: true",
			},
			{
				"unquoted string bytes: number",
				[]byte(`0`),
				"unsupported json token kind: number",
			},
			{
				"unquoted string bytes: truncated null",
				[]byte(`nul`),
				"failed to read token",
			},
			{
				"quoted string bytes: truncated string",
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f`),
				"invalid string",
			},
			{
				"quoted string bytes: empty",
				[]byte(`""`),
				"invalid string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.EthHash
						err := dec.unmarshal(tc.in, &n)
						require.ErrorContains(t, err, tc.want)
					})
				}
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
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewEthHash(ethcommon.Hash{}, false),
			},
			{
				"quoted hexadecimal string bytes: zero",
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000000"`),
				nullable.NewEthHash(ethcommon.Hash{}, true),
			},
			{
				"quoted hexadecimal string bytes: bitcoin genesis block",
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.EthHash
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.EthHash, n.EthHash)
					})
				}
			})
		}
	})
}

func TestEthHash_YAMLUnmarshaling(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"unquoted string bytes: sequence",
				[]byte(`[]`),
				"invalid node",
			},
			{
				"unquoted string bytes: mapping",
				[]byte(`{}`),
				"invalid node",
			},
			{
				"quoted string bytes: empty",
				[]byte(`""`),
				"invalid node",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthHash
				err := yaml.Unmarshal(tc.in, &n)
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
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewEthHash(ethcommon.Hash{}, false),
			},
			{
				"unquoted hexadecimal string bytes: zero",
				[]byte(`0x0000000000000000000000000000000000000000000000000000000000000000`),
				nullable.NewEthHash(ethcommon.Hash{}, true),
			},
			{
				"unquoted hexadecimal string bytes: bitcoin genesis block",
				[]byte(`0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f`),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
			{
				"quoted hexadecimal string bytes: zero",
				[]byte(`"0x0000000000000000000000000000000000000000000000000000000000000000"`),
				nullable.NewEthHash(ethcommon.Hash{}, true),
			},
			{
				"quoted hexadecimal string bytes: bitcoin genesis block",
				[]byte(`"0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"`),
				nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthHash
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.EthHash, n.EthHash)
			})
		}
	})

	t.Run("success: null node", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *yaml.Node
		}{
			{
				"no value with short tag",
				&yaml.Node{
					Kind: yaml.ScalarNode,
					Tag:  "!!null",
				},
			},
			{
				"no value with long tag",
				&yaml.Node{
					Kind: yaml.ScalarNode,
					Tag:  "tag:yaml.org,2002:null",
				},
			},
			{
				"null",
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "null",
				},
			},
			{
				"tilde",
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Value: "~",
				},
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewEthHash(ethcommon.HexToHash("0x000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, ethcommon.Hash{}, n.EthHash)
			})
		}
	})
}
