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

func TestEthAddress(t *testing.T) {
	var n nullable.EthAddress
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestEthAddress_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			want nullable.String
		}{
			{
				"invalid",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				nullable.NewString("", false),
			},
			{
				"valid: zero",
				nullable.NewEthAddress(ethcommon.Address{}, true),
				nullable.NewString("0x0000000000000000000000000000000000000000", true),
			},
			{
				"valid: vitalik.eth",
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
				"invalid",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewEthAddress(ethcommon.Address{}, true),
				ethhexutil.MustDecode("0x0000000000000000000000000000000000000000"),
			},
			{
				"valid: vitalik.eth",
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

func TestEthAddress_Scan(t *testing.T) {
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
				"bytes: zero",
				ethhexutil.MustDecode("0x0000000000000000000000000000000000000000"),
				nullable.NewEthAddress(ethcommon.Address{}, true),
			},
			{
				"bytes: vitalik.eth",
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

func TestEthAddress_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.EthAddress) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.EthAddress) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.EthAddress) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			want []byte
		}{
			{
				"invalid",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				[]byte(`null`),
			},
			{
				"valid: zero",
				nullable.NewEthAddress(ethcommon.Address{}, true),
				[]byte(`"0x0000000000000000000000000000000000000000"`),
			},
			{
				"valid: vitalik.eth",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
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

func TestEthAddress_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.EthAddress
			want []byte
		}{
			{
				"invalid",
				nullable.NewEthAddress(ethcommon.Address{}, false),
				[]byte("null\n"),
			},
			{
				"valid: zero",
				nullable.NewEthAddress(ethcommon.Address{}, true),
				[]byte("\"0x0000000000000000000000000000000000000000\"\n"),
			},
			{
				"valid: vitalik.eth",
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
				[]byte("\"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045\"\n"),
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

func TestEthAddress_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.EthAddress) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.EthAddress) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.EthAddress) error {
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
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045`),
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
						var n nullable.EthAddress
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
			want nullable.EthAddress
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewEthAddress(ethcommon.Address{}, false),
			},
			{
				"quoted hexadecimal string bytes: zero",
				[]byte(`"0x0000000000000000000000000000000000000000"`),
				nullable.NewEthAddress(ethcommon.Address{}, true),
			},
			{
				"quoted hexadecimal string bytes: vitalik.eth",
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.EthAddress
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.EthAddress, n.EthAddress)
					})
				}
			})
		}
	})
}

func TestEthAddress_YAMLUnmarshaling(t *testing.T) {
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
				var n nullable.EthAddress
				err := yaml.Unmarshal(tc.in, &n)
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
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewEthAddress(ethcommon.Address{}, false),
			},
			{
				"unquoted hexadecimal string bytes: zero",
				[]byte(`0x0000000000000000000000000000000000000000`),
				nullable.NewEthAddress(ethcommon.Address{}, true),
			},
			{
				"unquoted hexadecimal string bytes: vitalik.eth",
				[]byte(`0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045`),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
			{
				"quoted hexadecimal string bytes: zero",
				[]byte(`"0x0000000000000000000000000000000000000000"`),
				nullable.NewEthAddress(ethcommon.Address{}, true),
			},
			{
				"quoted hexadecimal string bytes: vitalik.eth",
				[]byte(`"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"`),
				nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.EthAddress
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.EthAddress, n.EthAddress)
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
				n := nullable.NewEthAddress(ethcommon.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, ethcommon.Address{}, n.EthAddress)
			})
		}
	})
}
