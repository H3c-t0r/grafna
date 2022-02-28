package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddlewareCSRF(t *testing.T) {
	tests := []struct {
		name        string
		cookieName  string
		method      string
		origin      string
		host        string
		defaultPort string
		code        int
	}{
		{
			name:        "mismatched origin and host is forbidden",
			cookieName:  "foo",
			method:      "GET",
			origin:      "http://notLocalhost",
			host:        "localhost",
			defaultPort: "80",
			code:        http.StatusForbidden,
		},
		{
			name:        "mismatched origin and host is NOT forbidden with a 'Safe Method'",
			cookieName:  "foo",
			method:      "TRACE",
			origin:      "http://notLocalhost",
			host:        "localhost",
			defaultPort: "80",
			code:        http.StatusOK,
		},
		{
			name:        "mismatched origin and host is NOT forbidden without a cookie",
			cookieName:  "",
			method:      "GET",
			origin:      "http://notLocalhost",
			host:        "localhost",
			defaultPort: "80",
			code:        http.StatusOK,
		},
		{
			name:        "malformed host is a bad request",
			cookieName:  "foo",
			method:      "GET",
			host:        "localhost:80:80",
			defaultPort: "80",
			code:        http.StatusBadRequest,
		},
		{
			name:        "should work without port",
			cookieName:  "foo",
			method:      "GET",
			host:        "localhost",
			origin:      "http://localhost",
			defaultPort: "80",
			code:        http.StatusOK,
		},
		{
			name:       "IPv6 host work with port",
			cookieName: "foo",
			method:     "GET",
			host:       "[::1]:3000",
			origin:     "http://[::1]:3000",
			code:       http.StatusOK,
		},
		{
			name:        "IPv6 host should get default port",
			cookieName:  "foo",
			method:      "GET",
			host:        "[::1]",
			origin:      "http://[::1]",
			defaultPort: "3000",
			code:        http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := csrfScenario(t, tt.cookieName, tt.method, tt.origin, tt.host, tt.defaultPort)
			require.Equal(t, tt.code, rr.Code)
		})
	}
}

func csrfScenario(t *testing.T, cookieName, method, origin, host, defaultPort string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{
		Name: cookieName,
	})

	// Note: Not sure where host header populates req.Host, or how that works.
	req.Host = host
	req.Header.Set("HOST", host)

	req.Header.Set("ORIGIN", origin)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})

	rr := httptest.NewRecorder()
	handler := CSRF(cookieName, defaultPort)(testHandler)
	handler.ServeHTTP(rr, req)
	return rr
}
