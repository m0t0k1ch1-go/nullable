package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/nullable"
	"github.com/m0t0k1ch1-go/nullable/internal/testutil"
)

func TestInt32MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			out  []byte
		}{
			{
				name: "null",
				in:   nullable.NewInt32(0, false),
				out:  []byte("null"),
			},
			{
				name: "not null",
				in:   nullable.NewInt32(1231006505, true),
				out:  []byte("1231006505"),
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

func TestInt32UnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Int32
		}{
			{
				name: "null",
				in:   []byte("null"),
				out:  nullable.NewInt32(0, false),
			},
			{
				name: "not null",
				in:   []byte("1231006505"),
				out:  nullable.NewInt32(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Int32
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out, n)
			})
		}
	})
}
