package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/m0t0k1ch1-go/timeutil/v5"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestTimestamp(t *testing.T) {
	var n nullable.Timestamp
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.MarshalerTo)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*yaml.Marshaler)(nil), &n)
	require.Implements(t, (*json.UnmarshalerFrom)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
	require.Implements(t, (*yaml.Unmarshaler)(nil), &n)
}

func TestTimestamp_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			want nullable.String
		}{
			{
				"invalid",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				nullable.NewString("", false),
			},
			{
				"valid: positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				nullable.NewString("1231006505", true),
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

func TestTimestamp_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			want driver.Value
		}{
			{
				"invalid",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				nil,
			},
			{
				"valid: positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				int64(1231006505),
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

func TestTimestamp_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"time",
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
				"unsupported source type: time.Time",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.Timestamp
		}{
			{
				"nil",
				nil,
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
			},
			{
				"int64: positive",
				int64(1231006505),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Timestamp.Unix(), n.Timestamp.Unix())
			})
		}
	})
}

func TestTimestamp_JSONMarshaling(t *testing.T) {
	encs := []struct {
		name    string
		marshal func(nullable.Timestamp) ([]byte, error)
	}{
		{
			"json.Marshal",
			func(n nullable.Timestamp) ([]byte, error) {
				return json.Marshal(n)
			},
		},
		{
			"MarshalJSON",
			func(n nullable.Timestamp) ([]byte, error) {
				return n.MarshalJSON()
			},
		},
	}

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			want []byte
		}{
			{
				"invalid",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				[]byte(`null`),
			},
			{
				"valid: positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				[]byte(`1231006505`),
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

func TestTimestamp_YAMLMarshaling(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			want []byte
		}{
			{
				"invalid",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				[]byte("null\n"),
			},
			{
				"valid: positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				[]byte("1231006505\n"),
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

func TestTimestamp_JSONUnmarshaling(t *testing.T) {
	decs := []struct {
		name      string
		unmarshal func([]byte, *nullable.Timestamp) error
	}{
		{
			"json.Unmarshal",
			func(b []byte, n *nullable.Timestamp) error {
				return json.Unmarshal(b, n)
			},
		},
		{
			"UnmarshalJSON",
			func(b []byte, n *nullable.Timestamp) error {
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
						var n nullable.Timestamp
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
			want nullable.Timestamp
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
			},
			{
				"unquoted decimal string bytes: positive",
				[]byte(`1231006505`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
			{
				"quoted decimal string bytes: positive",
				[]byte(`"1231006505"`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				for _, dec := range decs {
					t.Run(dec.name, func(t *testing.T) {
						var n nullable.Timestamp
						err := dec.unmarshal(tc.in, &n)
						require.NoError(t, err)
						require.Equal(t, tc.want.Valid, n.Valid)
						require.Equal(t, tc.want.Timestamp.Unix(), n.Timestamp.Unix())
					})
				}
			})
		}
	})
}

func TestTimestamp_YAMLUnmarshaling(t *testing.T) {
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
				var n nullable.Timestamp
				err := yaml.Unmarshal(tc.in, &n)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.Timestamp
		}{
			{
				"unquoted string bytes: null",
				[]byte(`null`),
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
			},
			{
				"unquoted decimal string bytes: positive",
				[]byte(`1231006505`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
			{
				"quoted decimal string bytes: positive",
				[]byte(`"1231006505"`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				err := yaml.Unmarshal(tc.in, &n)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Timestamp.Unix(), n.Timestamp.Unix())
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
				n := nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true)
				err := n.UnmarshalYAML(tc.in)
				require.NoError(t, err)
				require.False(t, n.Valid)
				require.Equal(t, timeutil.Timestamp{}.Unix(), n.Timestamp.Unix())
			})
		}
	})
}
