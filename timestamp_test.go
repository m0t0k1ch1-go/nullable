package nullable_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"
	"time"

	"github.com/m0t0k1ch1-go/timeutil/v4"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v2"
)

func TestTimestampNullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			out  nullable.String
		}{
			{
				"null",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				nullable.NewString("", false),
			},
			{
				"not null",
				nullable.NewTimestamp(timeutil.NewTimestamp(time.Unix(1231006505, 0)), true),
				nullable.NewString("1231006505", true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				require.Equal(t, tc.out, tc.in.NullableString())
			})
		}
	})
}

func TestTimestampValue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			out  driver.Value
		}{
			{
				"null",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				nil,
			},
			{
				"not null",
				nullable.NewTimestamp(timeutil.NewTimestamp(time.Unix(1231006505, 0)), true),
				int64(1231006505),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				v, err := tc.in.Value()
				if err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.out, v)
			})
		}
	})
}

func TestTimestampScan(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			out  nullable.Timestamp
		}{
			{
				"null",
				nil,
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
			},
			{
				"not null",
				int64(1231006505),
				nullable.NewTimestamp(timeutil.NewTimestamp(time.Unix(1231006505, 0)), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				if err := n.Scan(tc.in); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.out, n)
			})
		}
	})
}

func TestTimestampMarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			out  []byte
		}{
			{
				"null",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				[]byte("null"),
			},
			{
				"not null",
				nullable.NewTimestamp(timeutil.NewTimestamp(time.Unix(1231006505, 0)), true),
				[]byte("1231006505"),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := json.Marshal(tc.in)
				if err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.out, b)
			})
		}
	})
}

func TestTimestampUnmarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			out  nullable.Timestamp
		}{
			{
				"null",
				[]byte("null"),
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
			},
			{
				"not null",
				[]byte("1231006505"),
				nullable.NewTimestamp(timeutil.NewTimestamp(time.Unix(1231006505, 0)), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				if err := json.Unmarshal(tc.in, &n); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.out, n)
			})
		}
	})
}
