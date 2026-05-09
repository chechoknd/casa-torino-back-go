package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAppliesRouteSpecificLimit(t *testing.T) {
	handler := RateLimiter(
		RateLimit{Requests: 100, Window: time.Minute},
		WithRouteRateLimit("/auth/", RateLimit{Requests: 1, Window: time.Minute}),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	first := performRateLimitedRequest(handler, "/auth/login", "203.0.113.1")
	if first.Code != http.StatusOK {
		t.Fatalf("first auth request status = %d, want %d", first.Code, http.StatusOK)
	}

	second := performRateLimitedRequest(handler, "/auth/login", "203.0.113.1")
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second auth request status = %d, want %d", second.Code, http.StatusTooManyRequests)
	}

	otherRoute := performRateLimitedRequest(handler, "/products/", "203.0.113.1")
	if otherRoute.Code != http.StatusOK {
		t.Fatalf("default route status = %d, want %d", otherRoute.Code, http.StatusOK)
	}
}

func TestRateLimiterSeparatesClientsByForwardedIP(t *testing.T) {
	handler := RateLimiter(
		RateLimit{Requests: 100, Window: time.Minute},
		WithRouteRateLimit("/auth/", RateLimit{Requests: 1, Window: time.Minute}),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	firstClient := performRateLimitedRequest(handler, "/auth/login", "203.0.113.1")
	if firstClient.Code != http.StatusOK {
		t.Fatalf("first client status = %d, want %d", firstClient.Code, http.StatusOK)
	}

	secondClient := performRateLimitedRequest(handler, "/auth/login", "203.0.113.2")
	if secondClient.Code != http.StatusOK {
		t.Fatalf("second client status = %d, want %d", secondClient.Code, http.StatusOK)
	}
}

func performRateLimitedRequest(handler http.Handler, path, forwardedFor string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodPost, path, nil)
	request.Header.Set("X-Forwarded-For", forwardedFor)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	return recorder
}
