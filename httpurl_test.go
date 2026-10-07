package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"testing"

	"github.com/m0t0k1ch1-go/urlutil"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestHTTPURL(t *testing.T) {
	var n nullable.HTTPURL
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestHTTPURL_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want nullable.String
		}{
			{
				"invalid",
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
				nullable.NewString("", false),
			},
			{
				"valid: https",
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				nullable.NewString("https://m0t0k1ch1.com", true),
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

func TestHTTPURL_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want driver.Value
		}{
			{
				"invalid",
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
				nil,
			},
			{
				"valid: https",
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				"https://m0t0k1ch1.com",
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

func TestHTTPURL_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"bool",
				true,
				"unsupported source type: bool",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.HTTPURL
		}{
			{
				"nil",
				nil,
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
			},
			{
				"string: https",
				"https://m0t0k1ch1.com",
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.HTTPURL.String(), n.HTTPURL.String())
			})
		}
	})
}

func TestHTTPURL_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.HTTPURL) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.HTTPURL) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.HTTPURL) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want []byte
		}{
			{
				"invalid",
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
				[]byte(`null`),
			},
			{
				"valid: https",
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				[]byte(`"https://m0t0k1ch1.com"`),
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

func TestHTTPURL_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want []byte
		}{
			{
				"invalid",
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
				[]byte("null\n"),
			},
			{
				"valid: https",
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				[]byte("https://m0t0k1ch1.com\n"),
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

func TestHTTPURL_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.HTTPURL) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.HTTPURL) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.HTTPURL) error {
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
						var n nullable.HTTPURL
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
			want nullable.HTTPURL
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
			},
			{
				"quoted string bytes: https",
				[]byte(`"https://m0t0k1ch1.com"`),
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.HTTPURL
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.HTTPURL.String(), n.HTTPURL.String())
					})
				}
			})
		}
	})
}

func TestHTTPURL_YAMLUnmarshaling(t *testing.T) {
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
				var n nullable.HTTPURL
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.HTTPURL
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewHTTPURL(urlutil.HTTPURL{}, false),
			},
			{
				"unquoted string bytes: https",
				[]byte(`https://m0t0k1ch1.com`),
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
			},
			{
				"quoted string bytes: https",
				[]byte(`"https://m0t0k1ch1.com"`),
				nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.HTTPURL.String(), n.HTTPURL.String())
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
				n := nullable.NewHTTPURL(urlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, "", n.HTTPURL.String())
			})
		}
	})
}
