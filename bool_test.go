package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/nullable"
	"github.com/m0t0k1ch1-go/nullable/internal/testutil"
)

func TestBoolMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Bool
			out  []byte
		}{
			{
				name: "null",
				in:   nullable.NewBool(false, false),
				out:  []byte("null"),
			},
			{
				name: "not null",
				in:   nullable.NewBool(true, true),
				out:  []byte("true"),
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

func TestBoolUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Bool
		}{
			{
				name: "null",
				in:   []byte("null"),
				out:  nullable.NewBool(false, false),
			},
			{
				name: "not null",
				in:   []byte("true"),
				out:  nullable.NewBool(true, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Bool
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
