package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/coreutil"

	"github.com/m0t0k1ch1-go/nullable/v2"
	"github.com/m0t0k1ch1-go/nullable/v2/internal/testutil"
)

func TestNewInt32FromPtr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   *int32
			out  nullable.Int32
		}{
			{
				"nil",
				nil,
				nullable.NewInt32(0, false),
			},
			{
				"not nil",
				coreutil.Ptr(int32(1231006505)),
				nullable.NewInt32(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, nullable.NewInt32FromPtr(tc.in))
			})
		}
	})
}

func TestInt32Ptr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			out  *int32
		}{
			{
				"nil",
				nullable.NewInt32(0, false),
				nil,
			},
			{
				"not nil",
				nullable.NewInt32(1231006505, true),
				coreutil.Ptr(int32(1231006505)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				testutil.Equal(t, tc.out, tc.in.Ptr())
			})
		}
	})
}

func TestInt32MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Int32
			out  []byte
		}{
			{
				"null",
				nullable.NewInt32(0, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewInt32(1231006505, true),
				[]byte("1231006505"),
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
				"null",
				[]byte("null"),
				nullable.NewInt32(0, false),
			},
			{
				"not null",
				[]byte("1231006505"),
				nullable.NewInt32(1231006505, true),
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
