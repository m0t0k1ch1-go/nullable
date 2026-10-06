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

func TestInt64(t *testing.T) {
	var n nullable.Int64
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestNewInt64FromInt64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *int64
			want nullable.Int64
		}{
			{
				"nil",
				nil,
				nullable.NewInt64(0, false),
			},
			{
				"int64: zero",
				new(int64(0)),
				nullable.NewInt64(0, true),
			},
			{
				"int64: min",
				new(int64(math.MinInt64)),
				nullable.NewInt64(math.MinInt64, true),
			},
			{
				"int64: max",
				new(int64(math.MaxInt64)),
				nullable.NewInt64(math.MaxInt64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewInt64FromInt64Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int64, n.Int64)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		i := new(int64(1))
		n := nullable.NewInt64FromInt64Ptr(i)

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int64(1), n.Int64)
	})
}

func TestInt64_Int64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			want *int64
		}{
			{
				"invalid",
				nullable.NewInt64(0, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewInt64(0, true),
				new(int64(0)),
			},
			{
				"valid: min",
				nullable.NewInt64(math.MinInt64, true),
				new(int64(math.MinInt64)),
			},
			{
				"valid: max",
				nullable.NewInt64(math.MaxInt64, true),
				new(int64(math.MaxInt64)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				i := tc.in.Int64Ptr()
				require.Equal(t, tc.want, i)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		n := nullable.NewInt64(1, true)
		i := n.Int64Ptr()

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, int64(1), n.Int64)
	})
}

func TestInt64_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Int64) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Int64) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Int64) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			want []byte
		}{
			{
				"invalid",
				nullable.NewInt64(0, false),
				[]byte(`null`),
			},
			{
				"valid: zero",
				nullable.NewInt64(0, true),
				[]byte(`0`),
			},
			{
				"valid: min",
				nullable.NewInt64(math.MinInt64, true),
				[]byte(`-9223372036854775808`),
			},
			{
				"valid: max",
				nullable.NewInt64(math.MaxInt64, true),
				[]byte(`9223372036854775807`),
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

func TestInt64_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int64
			want []byte
		}{
			{
				"invalid",
				nullable.NewInt64(0, false),
				[]byte("null\n"),
			},
			{
				"valid: zero",
				nullable.NewInt64(0, true),
				[]byte("0\n"),
			},
			{
				"valid: min",
				nullable.NewInt64(math.MinInt64, true),
				[]byte("-9223372036854775808\n"),
			},
			{
				"valid: max",
				nullable.NewInt64(math.MaxInt64, true),
				[]byte("9223372036854775807\n"),
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

func TestInt64_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Int64) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Int64) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Int64) error {
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
				"invalid number",
			},
			{
				"unquoted decimal string bytes: exponential",
				[]byte(`0e0`),
				"invalid number",
			},
			{
				"unquoted decimal string bytes: min - 1",
				[]byte(`-9223372036854775809`),
				"invalid number",
			},
			{
				"unquoted decimal string bytes: max + 1",
				[]byte(`9223372036854775808`),
				"invalid number",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Int64
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
			want nullable.Int64
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewInt64(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewInt64(0, true),
			},
			{
				"unquoted decimal string bytes: min",
				[]byte(`-9223372036854775808`),
				nullable.NewInt64(math.MinInt64, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`9223372036854775807`),
				nullable.NewInt64(math.MaxInt64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Int64
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Int64, n.Int64)
					})
				}
			})
		}
	})
}

func TestInt64_YAMLUnmarshaling(t *testing.T) {
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
				"unquoted string bytes: boolean",
				[]byte(`true`),
				"invalid node",
			},
			{
				"quoted decimal string bytes: zero",
				[]byte(`"0"`),
				"invalid node",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int64
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Int64
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewInt64(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewInt64(0, true),
			},
			{
				"unquoted decimal string bytes: min",
				[]byte(`-9223372036854775808`),
				nullable.NewInt64(math.MinInt64, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`9223372036854775807`),
				nullable.NewInt64(math.MaxInt64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int64
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Int64, n.Int64)
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
				n := nullable.NewInt64(1, true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, int64(0), n.Int64)
			})
		}
	})
}
