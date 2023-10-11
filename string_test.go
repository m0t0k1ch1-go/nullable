package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/nullable"
	"github.com/m0t0k1ch1-go/nullable/internal/testutil"
)

func TestStringMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.String
			out  string
		}{
			{
				name: "null",
				in:   nullable.NewString("", false),
				out:  "null",
			},
			{
				name: "not null",
				in:   nullable.NewString("not null", true),
				out:  `"not null"`,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				if err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, string(b))
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
				var s nullable.String
				if err := json.Unmarshal(tc.in, &s); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, s)
			})
		}
	})
}
