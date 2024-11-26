package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/coreutil"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestNewStringFromStringPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *string
			out  nullable.String
		}{
			{
				"nil",
				nil,
				nullable.NewString("", false),
			},
			{
				"not nil",
				coreutil.Ptr("not nil"),
				nullable.NewString("not nil", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, nullable.NewStringFromStringPtr(tc.in))
			})
		}
	})
}

func TestStringStringPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			out  *string
		}{
			{
				"nil",
				nullable.NewString("", false),
				nil,
			},
			{
				"not nil",
				nullable.NewString("not nil", true),
				coreutil.Ptr("not nil"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, tc.in.StringPtr())
			})
		}
	})
}

func TestStringMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			out  []byte
		}{
			{
				"null",
				nullable.NewString("", false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewString("not null", true),
				[]byte(`"not null"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				require.Nil(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestStringUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.String
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewString("", false),
			},
			{
				"not null",
				[]byte(`"not null"`),
				nullable.NewString("not null", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				require.Nil(t, json.Unmarshal(tc.in, &n))

				require.Equal(t, tc.out, n)
			})
		}
	})
}

func TestStringMarshalYAML(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			out  any
		}{
			{
				"null",
				nullable.NewString("", false),
				[]byte("null\n"),
			},
			{
				"not null",
				nullable.NewString("not null", true),
				[]byte("not null\n"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := yaml.Marshal(tc.in)
				require.Nil(t, err)

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestStringUnmarshalYAML(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.String
		}{
			{
				"null",
				[]byte("null\n"),
				nullable.NewString("", false),
			},
			{
				"not null",
				[]byte("not null\n"),
				nullable.NewString("not null", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				require.Nil(t, yaml.Unmarshal(tc.in, &n))

				require.Equal(t, tc.out, n)
			})
		}
	})
}
