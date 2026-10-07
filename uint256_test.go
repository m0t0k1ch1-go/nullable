package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"testing"

	"github.com/m0t0k1ch1-go/bigutil/v3"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestUint256(t *testing.T) {
	var n nullable.Uint256
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestUint256_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			want nullable.String
		}{
			{
				"invalid",
				nullable.NewUint256(bigutil.Uint256{}, false),
				nullable.NewString("", false),
			},
			{
				"valid: one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				nullable.NewString("0x1", true),
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
				"invalid",
				nullable.NewUint256(bigutil.Uint256{}, false),
				nil,
			},
			{
				"valid: one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				[]byte{0x01},
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
				"unsupported source type: int64",
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
				"bytes: one",
				[]byte{0x01},
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
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

func TestUint256_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Uint256) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Uint256) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Uint256) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			want []byte
		}{
			{
				"invalid",
				nullable.NewUint256(bigutil.Uint256{}, false),
				[]byte(`null`),
			},
			{
				"valid: one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				[]byte(`"0x1"`),
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

func TestUint256_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint256
			want []byte
		}{
			{
				"invalid",
				nullable.NewUint256(bigutil.Uint256{}, false),
				[]byte("null\n"),
			},
			{
				"valid: one",
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
				[]byte("\"0x1\"\n"),
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

func TestUint256_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Uint256) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Uint256) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Uint256) error {
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
				"unquoted string bytes: truncated null",
				[]byte(`nul`),
				"failed to read token",
			},
			{
				"quoted string bytes: empty",
				[]byte(`""`),
				"invalid string: empty",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Uint256
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
			want nullable.Uint256
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewUint256(bigutil.Uint256{}, false),
			},
			{
				"quoted hexadecimal string bytes: one",
				[]byte(`"0x1"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"unquoted decimal string bytes: one",
				[]byte(`1`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Uint256
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Uint256.String(), n.Uint256.String())
					})
				}
			})
		}
	})
}

func TestUint256_YAMLUnmarshaling(t *testing.T) {
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
				"quoted string bytes: empty",
				[]byte(`""`),
				"invalid node",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				err := yaml.Unmarshal(tc.in, &n)
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
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewUint256(bigutil.Uint256{}, false),
			},
			{
				"unquoted hexadecimal string bytes: one",
				[]byte(`0x1`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"quoted hexadecimal string bytes: one",
				[]byte(`"0x1"`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
			{
				"unquoted decimal string bytes: one",
				[]byte(`1`),
				nullable.NewUint256(bigutil.NewUint256FromUint64(1), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint256
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint256.String(), n.Uint256.String())
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
				n := nullable.NewUint256(bigutil.NewUint256FromUint64(1), true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, "0x0", n.Uint256.String())
			})
		}
	})
}
