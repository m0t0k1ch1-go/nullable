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

func TestFloat64(t *testing.T) {
	var n nullable.Float64
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestNewFloat64FromFloat64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *float64
			want nullable.Float64
		}{
			{
				"nil",
				nil,
				nullable.NewFloat64(0, false),
			},
			{
				"float64: zero",
				new(float64(0)),
				nullable.NewFloat64(0, true),
			},
			{
				"float64: smallest non-zero",
				new(math.SmallestNonzeroFloat64),
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
			},
			{
				"float64: min",
				new(-math.MaxFloat64),
				nullable.NewFloat64(-math.MaxFloat64, true),
			},
			{
				"float64: max",
				new(math.MaxFloat64),
				nullable.NewFloat64(math.MaxFloat64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewFloat64FromFloat64Ptr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Float64, n.Float64)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		f := new(math.Phi)
		n := nullable.NewFloat64FromFloat64Ptr(f)

		*f = 0

		require.True(t, n.Valid)
		require.Equal(t, math.Phi, n.Float64)
	})
}

func TestFloat64_Float64Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Float64
			want *float64
		}{
			{
				"invalid",
				nullable.NewFloat64(0, false),
				nil,
			},
			{
				"valid: zero",
				nullable.NewFloat64(0, true),
				new(float64(0)),
			},
			{
				"valid: smallest non-zero",
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
				new(math.SmallestNonzeroFloat64),
			},
			{
				"valid: min",
				nullable.NewFloat64(-math.MaxFloat64, true),
				new(-math.MaxFloat64),
			},
			{
				"valid: max",
				nullable.NewFloat64(math.MaxFloat64, true),
				new(math.MaxFloat64),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				f := tc.in.Float64Ptr()
				require.Equal(t, tc.want, f)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		n := nullable.NewFloat64(math.Phi, true)
		f := n.Float64Ptr()

		*f = 0

		require.True(t, n.Valid)
		require.Equal(t, math.Phi, n.Float64)
	})
}

func TestFloat64_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Float64) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Float64) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Float64) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Float64
			want []byte
		}{
			{
				"invalid",
				nullable.NewFloat64(0, false),
				[]byte(`null`),
			},
			{
				"valid: zero",
				nullable.NewFloat64(0, true),
				[]byte(`0`),
			},
			{
				"valid: smallest non-zero",
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
				[]byte(`5e-324`),
			},
			{
				"valid: min",
				nullable.NewFloat64(-math.MaxFloat64, true),
				[]byte(`-1.7976931348623157e+308`),
			},
			{
				"valid: max",
				nullable.NewFloat64(math.MaxFloat64, true),
				[]byte(`1.7976931348623157e+308`),
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

func TestFloat64_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Float64
			want []byte
		}{
			{
				"invalid",
				nullable.NewFloat64(0, false),
				[]byte("null\n"),
			},
			{
				"valid: zero",
				nullable.NewFloat64(0, true),
				[]byte("0\n"),
			},
			{
				"valid: smallest non-zero",
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
				[]byte("5e-324\n"),
			},
			{
				"valid: min",
				nullable.NewFloat64(-math.MaxFloat64, true),
				[]byte("-1.7976931348623157e+308\n"),
			},
			{
				"valid: max",
				nullable.NewFloat64(math.MaxFloat64, true),
				[]byte("1.7976931348623157e+308\n"),
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

func TestFloat64_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Float64) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Float64) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Float64) error {
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
				"unquoted decimal string bytes: less than min",
				[]byte(`-1e+309`),
				"invalid float64",
			},
			{
				"unquoted decimal string bytes: greater than max",
				[]byte(`1e+309`),
				"invalid float64",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Float64
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
			want nullable.Float64
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewFloat64(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewFloat64(0, true),
			},
			{
				"unquoted decimal string bytes: smallest non-zero",
				[]byte(`5e-324`),
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
			},
			{
				"unquoted decimal string bytes: min",
				[]byte(`-1.7976931348623157e+308`),
				nullable.NewFloat64(-math.MaxFloat64, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`1.7976931348623157e+308`),
				nullable.NewFloat64(math.MaxFloat64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Float64
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Float64, n.Float64)
					})
				}
			})
		}
	})
}

func TestFloat64_YAMLUnmarshaling(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"unquoted string bytes: sequence",
				[]byte(`[]`),
				"invalid float64",
			},
			{
				"unquoted string bytes: mapping",
				[]byte(`{}`),
				"invalid float64",
			},
			{
				"unquoted string bytes: boolean",
				[]byte(`true`),
				"invalid float64",
			},
			{
				"quoted decimal string bytes: zero",
				[]byte(`"0"`),
				"invalid float64",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Float64
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Float64
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewFloat64(0, false),
			},
			{
				"unquoted decimal string bytes: zero",
				[]byte(`0`),
				nullable.NewFloat64(0, true),
			},
			{
				"unquoted decimal string bytes: smallest non-zero",
				[]byte(`5e-324`),
				nullable.NewFloat64(math.SmallestNonzeroFloat64, true),
			},
			{
				"unquoted decimal string bytes: min",
				[]byte(`-1.7976931348623157e+308`),
				nullable.NewFloat64(-math.MaxFloat64, true),
			},
			{
				"unquoted decimal string bytes: max",
				[]byte(`1.7976931348623157e+308`),
				nullable.NewFloat64(math.MaxFloat64, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Float64
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Float64, n.Float64)
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
				n := nullable.NewFloat64(math.Phi, true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, float64(0), n.Float64)
			})
		}
	})
}
