package nullable_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/m0t0k1ch1-go/nullable"
	"github.com/m0t0k1ch1-go/nullable/internal/testutil"
)

func TestStringMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			out  []byte
		}{
			{
				name: "null",
				in:   nullable.NewString("", false),
				out:  []byte("null"),
			},
			{
				name: "not null",
				in:   nullable.NewString("not null", true),
				out:  []byte(`"not null"`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, b)
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
				name: "null",
				in:   []byte("null"),
				out:  nullable.NewString("", false),
			},
			{
				name: "not null",
				in:   []byte(`"not null"`),
				out:  nullable.NewString("not null", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
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
				name: "null",
				in:   nullable.NewString("", false),
				out:  []byte("null\n"),
			},
			{
				name: "not null",
				in:   nullable.NewString("not null", true),
				out:  []byte("not null\n"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := yaml.Marshal(tc.in)
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, b)
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
				name: "null",
				in:   []byte("null\n"),
				out:  nullable.NewString("", false),
			},
			{
				name: "not null",
				in:   []byte("not null\n"),
				out:  nullable.NewString("not null", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.String
				if err := yaml.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
