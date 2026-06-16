package nullable_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/m0t0k1ch1-go/sqlutil/v3"
	"github.com/stretchr/testify/require"

	"github.com/m0t0k1ch1-go/nullable/v3"
)

func TestHTTPURL(t *testing.T) {
	var n nullable.HTTPURL
	require.Implements(t, (*driver.Valuer)(nil), &n)
	require.Implements(t, (*sql.Scanner)(nil), &n)
	require.Implements(t, (*json.Marshaler)(nil), &n)
	require.Implements(t, (*json.Unmarshaler)(nil), &n)
}

func TestHTTPURL_NullableString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want nullable.String
		}{
			{
				"null",
				nullable.NewHTTPURL(sqlutil.HTTPURL{}, false),
				nullable.NewString("", false),
			},
			{
				"http",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"), true),
				nullable.NewString("http://m0t0k1ch1.com", true),
			},
			{
				"https",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				nullable.NewString("https://m0t0k1ch1.com", true),
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

func TestHTTPURL_Value(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want driver.Value
		}{
			{
				"null",
				nullable.NewHTTPURL(sqlutil.HTTPURL{}, false),
				nil,
			},
			{
				"http",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"), true),
				"http://m0t0k1ch1.com",
			},
			{
				"https",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				"https://m0t0k1ch1.com",
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

func TestHTTPURL_Scan(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want string
		}{
			{
				"bool",
				true,
				"",
			},
			{
				"string: empty",
				"",
				"",
			},
			{
				"string: missing scheme",
				"m0t0k1ch1.com",
				"",
			},
			{
				"string: invalid host: empty",
				"://m0t0k1ch1.com",
				"",
			},
			{
				"string: invalid scheme: ftp",
				"ftp://m0t0k1ch1.com",
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := n.Scan(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   any
			want nullable.HTTPURL
		}{
			{
				"nil",
				nil,
				nullable.NewHTTPURL(sqlutil.HTTPURL{}, false),
			},
			{
				"string: http",
				"http://m0t0k1ch1.com",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"), true),
			},
			{
				"[]byte: https",
				[]byte("https://m0t0k1ch1.com"),
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := n.Scan(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.HTTPURL.String(), n.HTTPURL.String())
			})
		}
	})
}

func TestHTTPURL_MarshalJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   nullable.HTTPURL
			want []byte
		}{
			{
				"null",
				nullable.NewHTTPURL(sqlutil.HTTPURL{}, false),
				[]byte(`null`),
			},
			{
				"http",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"), true),
				[]byte(`"http://m0t0k1ch1.com"`),
			},
			{
				"https",
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
				[]byte(`"https://m0t0k1ch1.com"`),
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

func TestHTTPURL_UnmarshalJSON(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want string
		}{
			{
				"bool",
				[]byte(`true`),
				"",
			},
			{
				"string: empty",
				[]byte(`""`),
				"",
			},
			{
				"string: missing scheme",
				[]byte(`"m0t0k1ch1.com"`),
				"",
			},
			{
				"string: invalid host: empty",
				[]byte(`"://m0t0k1ch1.com"`),
				"",
			},
			{
				"string: invalid scheme: ftp",
				[]byte(`"ftp://m0t0k1ch1.com"`),
				"",
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := n.UnmarshalJSON(tc.in)
				require.ErrorContains(t, err, tc.want)
			})
		}
	})

	t.Run("success", func(t *testing.T) {
		tcs := []struct {
			name string
			in   []byte
			want nullable.HTTPURL
		}{
			{
				"null",
				[]byte(`null`),
				nullable.NewHTTPURL(sqlutil.HTTPURL{}, false),
			},
			{
				"http",
				[]byte(`"http://m0t0k1ch1.com"`),
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("http://m0t0k1ch1.com"), true),
			},
			{
				"https",
				[]byte(`"https://m0t0k1ch1.com"`),
				nullable.NewHTTPURL(sqlutil.MustNewHTTPURLFromString("https://m0t0k1ch1.com"), true),
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				var n nullable.HTTPURL
				err := n.UnmarshalJSON(tc.in)
				require.NoError(t, err)
				require.Equal(t, tc.want.Valid, n.Valid)
				require.Equal(t, tc.want.HTTPURL.String(), n.HTTPURL.String())
			})
		}
	})
}
