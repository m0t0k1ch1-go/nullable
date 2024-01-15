package nullable_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/m0t0k1ch1-go/nullable"
	"github.com/m0t0k1ch1-go/nullable/internal/testutil"
)

func TestBigIntNullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.BigInt
			out  nullable.String
		}{
			{
				"null",
				nullable.NewBigInt(nil, false),
				nullable.NewString("", false),
			},
			{
				"not null",
				nullable.NewBigInt(big.NewInt(1231006505), true),
				nullable.NewString("1231006505", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, tc.in.NullableString())
			})
		}
	})
}

func TestBigIntMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.BigInt
			out  []byte
		}{
			{
				name: "null",
				in:   nullable.NewBigInt(nil, false),
				out:  []byte("null"),
			},
			{
				name: "not null",
				in:   nullable.NewBigInt(big.NewInt(1231006505), true),
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

func TestBigIntUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.BigInt
		}{
			{
				name: "null",
				in:   []byte("null"),
				out:  nullable.NewBigInt(nil, false),
			},
			{
				name: "not null",
				in:   []byte("1231006505"),
				out:  nullable.NewBigInt(big.NewInt(1231006505), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.BigInt
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				testutil.Equal(t, tc.out.Valid, n.Valid)
				testutil.Equal(t, n.BigInt.Cmp(tc.out.BigInt), 0)
			})
		}
	})
}
