package nullable_test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/m0t0k1ch1-go/timeutil/v5"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestTimestamp(t *testing.T) {
	var n nullable.Timestamp
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestTimestamp_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			want nullable.String
		}{
			{
				"null",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				nullable.NewString("", false),
			},
			{
				"zero",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(0), true),
				nullable.NewString("0", true),
			},
			{
				"positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				nullable.NewString("1231006505", true),
			},
			{
				"negative",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(-1231006505), true),
				nullable.NewString("-1231006505", true),
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
				"null",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				nil,
			},
			{
				"zero",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(0), true),
				int64(0),
			},
			{
				"positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				int64(1231006505),
			},
			{
				"negative",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(-1231006505), true),
				int64(-1231006505),
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
				"time.Time",
				time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
				"",
			},
			{
				"uint64: exceeds int64 range",
				uint64(math.MaxInt64) + 1,
				"",
			},
			{
				"[]byte: empty",
				[]byte{},
				"",
			},
			{
				"[]byte: exceeds 256 bits",
				append([]byte{0x01}, bytes.Repeat([]byte{0x00}, 32)...),
				"",
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
				"int64: zero",
				int64(0),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(0), true),
			},
			{
				"int64: positive",
				int64(1231006505),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
			{
				"int64: negative",
				int64(-1231006505),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(-1231006505), true),
			},
			{
				"uint64",
				uint64(1231006505),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
			{
				"[]byte",
				[]byte("1231006505"),
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

func TestTimestamp_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.Timestamp
			want []byte
		}{
			{
				"null",
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
				[]byte(`null`),
			},
			{
				"zero",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(0), true),
				[]byte(`0`),
			},
			{
				"positive",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
				[]byte(`1231006505`),
			},
			{
				"negative",
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(-1231006505), true),
				[]byte(`-1231006505`),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				b, err := tc.in.MarshalJSON()
				require.NoError(t, err)
				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"empty",
				[]byte{},
				"",
			},
			{
				"number: exceeds int64 range",
				[]byte(`9223372036854775808`),
				"",
			},
			{
				"number: fractional",
				[]byte(`1231006505.0`),
				"",
			},
			{
				"number: exponential",
				[]byte(`1231006505e0`),
				"",
			},
			{
				"string: empty",
				[]byte(`""`),
				"",
			},
			{
				"string: zero",
				[]byte(`"0"`),
				"",
			},
			{
				"string: positive decimal",
				[]byte(`"1231006505"`),
				"",
			},
			{
				"string: negative decimal",
				[]byte(`"-1231006505"`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				err := n.UnmarshalJSON(tc.in)
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
				"null",
				[]byte(`null`),
				nullable.NewTimestamp(timeutil.Timestamp{}, false),
			},
			{
				"number: zero",
				[]byte(`0`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(0), true),
			},
			{
				"number: positive",
				[]byte(`1231006505`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(1231006505), true),
			},
			{
				"number: negative",
				[]byte(`-1231006505`),
				nullable.NewTimestamp(timeutil.NewTimestampFromUnix(-1231006505), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.Timestamp
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.Timestamp.Unix(), n.Timestamp.Unix())
			})
		}
	})
}
