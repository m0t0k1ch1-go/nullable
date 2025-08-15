package nullable_test

import (
	"encoding/json"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestNewInt32FromInt32Ptr(t *testing.T) {
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
				lo.ToPtr(int32(1231006505)),
				nullable.NewInt32(1231006505, true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, nullable.NewInt32FromInt32Ptr(tc.in))
			})
		}
	})
}

func TestInt32Int32Ptr(t *testing.T) {
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
				lo.ToPtr(int32(1231006505)),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, tc.in.Int32Ptr())
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
				require.Nil(t, err)

				require.Equal(t, tc.out, b)
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
				require.Nil(t, json.Unmarshal(tc.in, &n))

				require.Equal(t, tc.out, n)
			})
		}
	})
}
