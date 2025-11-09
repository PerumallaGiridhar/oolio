package routes

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMemUsage_ReturnsStats(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rr := httptest.NewRecorder()

	MemUsage(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if ct == "" || ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	for _, key := range []string{"Alloc", "TotalAlloc", "Sys", "NumGC"} {
		if _, ok := body[key]; !ok {
			t.Errorf("expected key %q in response body", key)
		}
	}
}

func TestStatsEndpoint(t *testing.T) {
	r := NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/stats", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	type stats struct {
		Alloc      string
		TotalAlloc string
		Sys        string
		NumGC      string
	}
	var s stats
	if err := json.Unmarshal(body, &s); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestNewRouter_HeartbeatLive(t *testing.T) {
	r := NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("expected status 200 or 204 for /live, got %d", rr.Code)
	}
}

func TestNewRouter_CORSHeaders(t *testing.T) {
	r := NewRouter()

	req := httptest.NewRequest(http.MethodOptions, "/stats", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent && rr.Code != http.StatusOK {
		t.Fatalf("expected CORS preflight status 200 or 204, got %d", rr.Code)
	}

	allowOrigin := rr.Header().Get("Access-Control-Allow-Origin")
	if allowOrigin == "" {
		t.Fatalf("expected Access-Control-Allow-Origin header to be set")
	}
}

func TestNewRouter_APIProductRouteExists(t *testing.T) {
	r := NewRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/product", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent && rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 or 204 for OPTIONS /api/product, got %d", rr.Code)
	}

	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("expected Access-Control-Allow-Origin header in OPTIONS /api/product response")
	}
}

func TestNewRouter_APIProductIdRouteExists(t *testing.T) {
	r := NewRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/product/1", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent && rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 or 204 for OPTIONS /api/product/1, got %d", rr.Code)
	}

	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("expected Access-Control-Allow-Origin header in OPTIONS /api/product/api response")
	}
}

func TestNewRouter_APICreateOrderRouteExists(t *testing.T) {
	r := NewRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/order", nil)
	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent && rr.Code != http.StatusOK {
		t.Fatalf("expected status 200 or 204 for OPTIONS /api/order/1, got %d", rr.Code)
	}

	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Fatalf("expected Access-Control-Allow-Origin header in OPTIONS /api/order response")
	}
}
