package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestInt32(t *testing.T) {
	var n nullable.Int32
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestNewInt32FromInt32Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *int32
			want nullable.Int32
		}{
			{
				"nil",
				nil,
				nullable.NewInt32(0, false),
			},
			{
				"int32: zero",
				new(int32(0)),
				nullable.NewInt32(0, true),
			},
			{
				"int32: min",
				new(int32(math.MinInt32)),
				nullable.NewInt32(math.MinInt32, true),
			},
			{
				"int32: max",
				new(int32(math.MaxInt32)),
				nullable.NewInt32(math.MaxInt32, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewInt32FromInt32Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int32, n.Int32)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		i := new(int32(1))
		n := nullable.NewInt32FromInt32Ptr(i)

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int32(1), n.Int32)
	})
}

func TestInt32_Int32Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			want *int32
		}{
			{
				"invalid",
				nullable.NewInt32(0, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewInt32(0, true),
				new(int32(0)),
			},
			{
				"valid: min",
				nullable.NewInt32(math.MinInt32, true),
				new(int32(math.MinInt32)),
			},
			{
				"valid: max",
				nullable.NewInt32(math.MaxInt32, true),
				new(int32(math.MaxInt32)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				i := tc.in.Int32Ptr()
				require.Equal(t, tc.want, i)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		n := nullable.NewInt32(1, true)
		i := n.Int32Ptr()

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int32(1), n.Int32)
	})
}

func TestInt32_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Int32) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Int32) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Int32) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			want []byte
		}{
			{
				"invalid",
				nullable.NewInt32(0, false),
				[]byte(`null`),
			},
			{
				"valid: zero",
				nullable.NewInt32(0, true),
				[]byte(`0`),
			},
			{
				"valid: min",
				nullable.NewInt32(math.MinInt32, true),
				[]byte(`-2147483648`),
			},
			{
				"valid: max",
				nullable.NewInt32(math.MaxInt32, true),
				[]byte(`2147483647`),
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

func TestInt32_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			want []byte
		}{
			{
				"invalid",
				nullable.NewInt32(0, false),
				[]byte("null\n"),
			},
			{
				"valid: zero",
				nullable.NewInt32(0, true),
				[]byte("0\n"),
			},
			{
				"valid: min",
				nullable.NewInt32(math.MinInt32, true),
				[]byte("-2147483648\n"),
			},
			{
				"valid: max",
				nullable.NewInt32(math.MaxInt32, true),
				[]byte("2147483647\n"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := yaml.Marshal(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestInt32_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Int32) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Int32) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Int32) error {
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
				"quoted decimal string bytes: zero",
				[]byte(`"0"`),
				"unsupported json token kind: string",
			},
			{
				"unquoted string bytes: truncated null",
				[]byte(`nul`),
				"failed to read token",
			},
			{
				"unquoted decimal string bytes: fractional",
				[]byte(`0.0`),
				"invalid int32",
			},
			{
				"unquoted decimal string bytes: exponential",
				[]byte(`0e0`),
				"invalid int32",
			},
			{
				"unquoted decimal string bytes: min - 1",
				[]byte(`-2147483649`),
				"invalid int32",
			},
			{
				"unquoted decimal string bytes: max + 1",
				[]byte(`2147483648`),
				"invalid int32",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Int32
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
			want nullable.Int32
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewInt32(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewInt32(0, true),
			},
			{
				"unquoted decimal string bytes: min",
				[]byte(`-2147483648`),
				nullable.NewInt32(math.MinInt32, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`2147483647`),
				nullable.NewInt32(math.MaxInt32, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Int32
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Int32, n.Int32)
					})
				}
			})
		}
	})
}

func TestInt32_YAMLUnmarshaling(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"unquoted string bytes: sequence",
				[]byte(`[]`),
				"invalid int32",
			},
			{
				"unquoted string bytes: mapping",
				[]byte(`{}`),
				"invalid int32",
			},
			{
				"unquoted string bytes: boolean",
				[]byte(`true`),
				"invalid int32",
			},
			{
				"quoted decimal string bytes: zero",
				[]byte(`"0"`),
				"invalid int32",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int32
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Int32
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewInt32(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewInt32(0, true),
			},
			{
				"unquoted decimal string bytes: min",
				[]byte(`-2147483648`),
				nullable.NewInt32(math.MinInt32, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`2147483647`),
				nullable.NewInt32(math.MaxInt32, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int32
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int32, n.Int32)
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
				n := nullable.NewInt32(1, true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, int32(0), n.Int32)
			})
		}
	})
}
