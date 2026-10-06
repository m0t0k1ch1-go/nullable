package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestBool(t *testing.T) {
	var n nullable.Bool
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestNewBoolFromBoolPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *bool
			want nullable.Bool
		}{
			{
				"nil",
				nil,
				nullable.NewBool(false, false),
			},
			{
				"bool: true",
				new(true),
				nullable.NewBool(true, true),
			},
			{
				"bool: false",
				new(false),
				nullable.NewBool(false, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewBoolFromBoolPtr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Bool, n.Bool)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		b := new(true)
		n := nullable.NewBoolFromBoolPtr(b)

		*b = false

		require.True(t, n.Valid)
		require.True(t, n.Bool)
	})
}

func TestBool_BoolPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			want *bool
		}{
			{
				"invalid",
				nullable.NewBool(false, false),
				nil,
			},
			{
				"valid: true",
				nullable.NewBool(true, true),
				new(true),
			},
			{
				"valid: false",
				nullable.NewBool(false, true),
				new(false),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b := tc.in.BoolPtr()
				require.Equal(t, tc.want, b)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		n := nullable.NewBool(true, true)
		b := n.BoolPtr()

		*b = false

		require.True(t, n.Valid)
		require.True(t, n.Bool)
	})
}

func TestBool_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Bool) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Bool) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Bool) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			want []byte
		}{
			{
				"invalid",
				nullable.NewBool(false, false),
				[]byte(`null`),
			},
			{
				"valid: true",
				nullable.NewBool(true, true),
				[]byte(`true`),
			},
			{
				"valid: false",
				nullable.NewBool(false, true),
				[]byte(`false`),
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

func TestBool_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			want []byte
		}{
			{
				"invalid",
				nullable.NewBool(false, false),
				[]byte("null\n"),
			},
			{
				"valid: true",
				nullable.NewBool(true, true),
				[]byte("true\n"),
			},
			{
				"valid: false",
				nullable.NewBool(false, true),
				[]byte("false\n"),
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

func TestBool_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Bool) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Bool) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Bool) error {
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
				"unquoted string bytes: number",
				[]byte(`0`),
				"unsupported json token kind: number",
			},
			{
				"quoted string bytes: boolean",
				[]byte(`"true"`),
				"unsupported json token kind: string",
			},
			{
				"unquoted string bytes: truncated null",
				[]byte(`nul`),
				"failed to read token",
			},
			{
				"unquoted string bytes: truncated boolean",
				[]byte(`tru`),
				"invalid bool",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Bool
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
			want nullable.Bool
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewBool(false, false),
			},
			{
				"unquoted string bytes: true",
				[]byte(`true`),
				nullable.NewBool(true, true),
			},
			{
				"unquoted string bytes: false",
				[]byte(`false`),
				nullable.NewBool(false, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Bool
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Bool, n.Bool)
					})
				}
			})
		}
	})
}

func TestBool_YAMLUnmarshaling(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"unquoted string bytes: sequence",
				[]byte(`[]`),
				"invalid bool",
			},
			{
				"unquoted string bytes: mapping",
				[]byte(`{}`),
				"invalid bool",
			},
			{
				"unquoted string bytes: number",
				[]byte(`0`),
				"invalid bool",
			},
			{
				"quoted string bytes: boolean",
				[]byte(`"true"`),
				"invalid bool",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Bool
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Bool
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewBool(false, false),
			},
			{
				"unquoted string bytes: true",
				[]byte(`true`),
				nullable.NewBool(true, true),
			},
			{
				"unquoted string bytes: false",
				[]byte(`false`),
				nullable.NewBool(false, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Bool
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Bool, n.Bool)
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
				n := nullable.NewBool(true, true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.False(t, n.Bool)
			})
		}
	})
}
