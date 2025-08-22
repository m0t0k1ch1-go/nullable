package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

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
				"empty",
				ptr(""),
				nullable.NewString("", true),
			},
			{
				"non-empty",
				ptr("non-empty"),
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
}

func TestString_StringPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			want *string
		}{
			{
				"null",
				nullable.NewString("", false),
				nil,
			},
			{
				"empty",
				nullable.NewString("", true),
				ptr(""),
			},
			{
				"non-empty",
				nullable.NewString("non-empty", true),
				ptr("non-empty"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				n := tc.in
				p := n.StringPtr()
				require.Equal(t, tc.want, p)

				if p != nil {
					*p = tc.in.String + " modified"

					require.Equal(t, tc.in.Valid, n.Valid)
					require.Equal(t, tc.in.String, n.String)
				}
			})
		}
	})
}

func TestString_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			want []byte
		}{
			{
				"null",
				nullable.NewString("", false),
				[]byte(`null`),
			},
			{
				"empty",
				nullable.NewString("", true),
				[]byte(`""`),
			},
			{
				"non-empty",
				nullable.NewString("non-empty", true),
				[]byte(`"non-empty"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestString_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"boolean",
				[]byte(`true`),
				"",
			},
			{
				"number",
				[]byte(`0`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				err := n.UnmarshalJSON(tc.in)
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
				"null",
				[]byte(`null`),
				nullable.NewString("", false),
			},
			{
				"string: empty",
				[]byte(`""`),
				nullable.NewString("", true),
			},
			{
				"string: non-empty",
				[]byte(`"non-empty"`),
				nullable.NewString("non-empty", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.String, n.String)
			})
		}
	})
}

func TestString_MarshalYAML(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			want any
		}{
			{
				"null",
				nullable.NewString("", false),
				nil,
			},
			{
				"empty",
				nullable.NewString("", true),
				"",
			},
			{
				"non-empty",
				nullable.NewString("non-empty", true),
				"non-empty",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.MarshalYAML()
				require.NoError(t, err)
				require.Equal(t, tc.want, v)
			})
		}
	})
}

func TestString_UnmarshalYAML(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *yaml.Node
			want string
		}{
			{
				"sequence",
				&yaml.Node{
					Kind: yaml.SequenceNode,
					Tag:  "!!seq",
				},
				"",
			},
			{
				"mapping",
				&yaml.Node{
					Kind: yaml.MappingNode,
					Tag:  "!!map",
				},
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				err := n.UnmarshalYAML(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *yaml.Node
			want nullable.String
		}{
			{
				"null",
				&yaml.Node{
					Kind: yaml.ScalarNode,
					Tag:  "!!null",
				},
				nullable.NewString("", false),
			},
			{
				"string: empty",
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: "",
				},
				nullable.NewString("", true),
			},
			{
				"string: non-empty",
				&yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: "non-empty",
				},
				nullable.NewString("non-empty", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.String, n.String)
			})
		}
	})
}
