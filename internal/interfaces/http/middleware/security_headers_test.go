package middleware

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Forwarded-Proto", "https")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	assertHeader(t, recorder, "X-Frame-Options", "DENY")
	assertHeader(t, recorder, "X-Content-Type-Options", "nosniff")
	assertHeader(t, recorder, "Referrer-Policy", "strict-origin-when-cross-origin")
	assertHeader(t, recorder, "Strict-Transport-Security", "max-age=63072000; includeSubDomains")
}

func TestSecurityHeadersOnlySetsHSTSForHTTPS(t *testing.T) {
	handler := SecurityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	httpRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	httpRecorder := httptest.NewRecorder()
	handler.ServeHTTP(httpRecorder, httpRequest)

	if got := httpRecorder.Header().Get("Strict-Transport-Security"); got != "" {
		t.Fatalf("Strict-Transport-Security = %q, want empty for HTTP", got)
	}

	httpsRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	httpsRequest.TLS = &tls.ConnectionState{}
	httpsRecorder := httptest.NewRecorder()
	handler.ServeHTTP(httpsRecorder, httpsRequest)

	assertHeader(t, httpsRecorder, "Strict-Transport-Security", "max-age=63072000; includeSubDomains")
}

func assertHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, want string) {
	t.Helper()

	if got := recorder.Header().Get(key); got != want {
		t.Fatalf("%s = %q, want %q", key, got, want)
	}
}
