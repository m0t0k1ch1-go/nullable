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

func TestString(t *testing.T) {
	var n nullable.String
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestNewStringFromStringPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *string
			want nullable.String
		}{
			{
				"nil",
				nil,
				nullable.NewString("", false),
			},
			{
				"string: empty",
				new(""),
				nullable.NewString("", true),
			},
			{
				"string: non-empty",
				new("non-empty"),
				nullable.NewString("non-empty", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := nullable.NewStringFromStringPtr(tc.in)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.String, n.String)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		s := new("non-empty")
		n := nullable.NewStringFromStringPtr(s)

		*s = ""

		require.True(t, n.Valid)
		require.Equal(t, "non-empty", n.String)
	})
}

func TestString_StringPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			want *string
		}{
			{
				"invalid",
				nullable.NewString("", false),
				nil,
			},
			{
				"valid: empty",
				nullable.NewString("", true),
				new(""),
			},
			{
				"valid: non-empty",
				nullable.NewString("non-empty", true),
				new("non-empty"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				s := tc.in.StringPtr()
				require.Equal(t, tc.want, s)
			})
		}
	})

	t.Run("success: no aliasing", func(t *testing.T) {
		n := nullable.NewString("non-empty", true)
		s := n.StringPtr()

		*s = ""

		require.True(t, n.Valid)
		require.Equal(t, "non-empty", n.String)
	})
}

func TestString_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.String) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.String) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.String) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			want []byte
		}{
			{
				"invalid",
				nullable.NewString("", false),
				[]byte(`null`),
			},
			{
				"valid: empty",
				nullable.NewString("", true),
				[]byte(`""`),
			},
			{
				"valid: non-empty",
				nullable.NewString("non-empty", true),
				[]byte(`"non-empty"`),
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

func TestString_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			want []byte
		}{
			{
				"invalid",
				nullable.NewString("", false),
				[]byte("null\n"),
			},
			{
				"valid: empty",
				nullable.NewString("", true),
				[]byte(`""` + "\n"),
			},
			{
				"valid: non-empty",
				nullable.NewString("non-empty", true),
				[]byte("non-empty\n"),
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

func TestString_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.String) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.String) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.String) error {
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
				[]byte(`"non-empty`),
				"invalid string",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.String
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
			want nullable.String
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewString("", false),
			},
			{
				"quoted string bytes: empty",
				[]byte(`""`),
				nullable.NewString("", true),
			},
			{
				"quoted string bytes: non-empty",
				[]byte(`"non-empty"`),
				nullable.NewString("non-empty", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.String
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.String, n.String)
					})
				}
			})
		}
	})
}

func TestString_YAMLUnmarshaling(t *testing.T) {
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
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.String
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewString("", false),
			},
			{
				"unquoted string bytes: boolean",
				[]byte(`true`),
				nullable.NewString("true", true),
			},
			{
				"unquoted string bytes: number",
				[]byte(`0`),
				nullable.NewString("0", true),
			},
			{
				"quoted string bytes: empty",
				[]byte(`""`),
				nullable.NewString("", true),
			},
			{
				"quoted string bytes: non-empty",
				[]byte(`"non-empty"`),
				nullable.NewString("non-empty", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.String, n.String)
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
				n := nullable.NewString("non-empty", true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, "", n.String)
			})
		}
	})
}
