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

func TestUint64(t *testing.T) {
	var n nullable.Uint64
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestNewUint64FromUint64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *uint64
			want nullable.Uint64
		}{
			{
				"nil",
				nil,
				nullable.NewUint64(0, false),
			},
			{
				"uint64: zero",
				new(uint64(0)),
				nullable.NewUint64(0, true),
			},
			{
				"uint64: one",
				new(uint64(1)),
				nullable.NewUint64(1, true),
			},
			{
				"uint64: max",
				new(uint64(math.MaxUint64)),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewUint64FromUint64Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint64, n.Uint64)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		i := new(uint64(1))
		n := nullable.NewUint64FromUint64Ptr(i)

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, uint64(1), n.Uint64)
	})
}

func TestUint64_Uint64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want *uint64
		}{
			{
				"invalid",
				nullable.NewUint64(0, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewUint64(0, true),
				new(uint64(0)),
			},
			{
				"valid: one",
				nullable.NewUint64(1, true),
				new(uint64(1)),
			},
			{
				"valid: max",
				nullable.NewUint64(math.MaxUint64, true),
				new(uint64(math.MaxUint64)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				i := tc.in.Uint64Ptr()
				require.Equal(t, tc.want, i)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		n := nullable.NewUint64(1, true)
		i := n.Uint64Ptr()

		*i = 0

		require.True(t, n.Valid)
		require.Equal(t, uint64(1), n.Uint64)
	})
}

func TestUint64_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want driver.Value
		}{
			{
				"invalid",
				nullable.NewUint64(0, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewUint64(0, true),
				uint64(0),
			},
			{
				"valid: one",
				nullable.NewUint64(1, true),
				uint64(1),
			},
			{
				"valid: max",
				nullable.NewUint64(math.MaxUint64, true),
				uint64(math.MaxUint64),
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

func TestUint64_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"float64",
				float64(0),
				"unsupported source type: float64",
			},
			{
				"int64: negative",
				int64(-1),
				"invalid int64 source: negative",
			},
			{
				"bytes: empty",
				[]byte{},
				"invalid bytes source: empty",
			},
			{
				"bytes: invalid",
				[]byte("invalid"),
				"invalid bytes source",
			},
			{
				"decimal string bytes: signed negative",
				[]byte("-1"),
				"invalid bytes source",
			},
			{
				"decimal string bytes: max + 1",
				[]byte("18446744073709551616"),
				"invalid bytes source",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.Uint64
		}{
			{
				"nil",
				nil,
				nullable.NewUint64(0, false),
			},
			{
				"int64: zero",
				int64(0),
				nullable.NewUint64(0, true),
			},
			{
				"int64: one",
				int64(1),
				nullable.NewUint64(1, true),
			},
			{
				"int64: max",
				int64(math.MaxInt64),
				nullable.NewUint64(math.MaxInt64, true),
			},
			{
				"uint64: zero",
				uint64(0),
				nullable.NewUint64(0, true),
			},
			{
				"uint64: one",
				uint64(1),
				nullable.NewUint64(1, true),
			},
			{
				"uint64: max",
				uint64(math.MaxUint64),
				nullable.NewUint64(math.MaxUint64, true),
			},
			{
				"decimal string bytes: zero",
				[]byte("0"),
				nullable.NewUint64(0, true),
			},
			{
				"decimal string bytes: one",
				[]byte("1"),
				nullable.NewUint64(1, true),
			},
			{
				"decimal string bytes: max",
				[]byte("18446744073709551615"),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint64, n.Uint64)
			})
		}
	})
}

func TestUint64_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Uint64) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Uint64) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Uint64) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want []byte
		}{
			{
				"invalid",
				nullable.NewUint64(0, false),
				[]byte(`null`),
			},
			{
				"valid: zero",
				nullable.NewUint64(0, true),
				[]byte(`0`),
			},
			{
				"valid: one",
				nullable.NewUint64(1, true),
				[]byte(`1`),
			},
			{
				"valid: max",
				nullable.NewUint64(math.MaxUint64, true),
				[]byte(`18446744073709551615`),
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

func TestUint64_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Uint64
			want []byte
		}{
			{
				"invalid",
				nullable.NewUint64(0, false),
				[]byte("null\n"),
			},
			{
				"valid: zero",
				nullable.NewUint64(0, true),
				[]byte("0\n"),
			},
			{
				"valid: one",
				nullable.NewUint64(1, true),
				[]byte("1\n"),
			},
			{
				"valid: max",
				nullable.NewUint64(math.MaxUint64, true),
				[]byte("18446744073709551615\n"),
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

func TestUint64_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Uint64) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Uint64) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Uint64) error {
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
				"unquoted decimal string bytes: signed negative",
				[]byte(`-1`),
				"invalid number",
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
				"unquoted decimal string bytes: max + 1",
				[]byte(`18446744073709551616`),
				"invalid number",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Uint64
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
			want nullable.Uint64
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewUint64(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewUint64(0, true),
			},
			{
				"unquoted decimal string bytes: one",
				[]byte(`1`),
				nullable.NewUint64(1, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`18446744073709551615`),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Uint64
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Uint64, n.Uint64)
					})
				}
			})
		}
	})
}

func TestUint64_YAMLUnmarshaling(t *testing.T) {
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
				var n nullable.Uint64
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Uint64
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewUint64(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewUint64(0, true),
			},
			{
				"unquoted decimal string bytes: one",
				[]byte(`1`),
				nullable.NewUint64(1, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`18446744073709551615`),
				nullable.NewUint64(math.MaxUint64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Uint64
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Uint64, n.Uint64)
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
				n := nullable.NewUint64(1, true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, uint64(0), n.Uint64)
			})
		}
	})
}
