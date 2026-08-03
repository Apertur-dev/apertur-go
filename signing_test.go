package apertur

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Known-answer vectors, generated against the node reference implementation
// (packages/client-node/src/signature.ts signRequest).
func TestSignRequest_KnownAnswer_WithBody(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	sig, ts := SignRequest(secret, "POST", "/api/v1/upload-sessions", []byte(`{"a":1}`), 1800000000)

	wantSig := "sha256=d5bf88c946aa6cf4397749eacd05cf058e6828878f9b1c84830cb7c07a234d3e"
	wantTs := "1800000000"

	if sig != wantSig {
		t.Errorf("sig header = %q, want %q", sig, wantSig)
	}
	if ts != wantTs {
		t.Errorf("ts header = %q, want %q", ts, wantTs)
	}
}

func TestSignRequest_KnownAnswer_EmptyBody(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	sig, ts := SignRequest(secret, "GET", "/api/v1/uploads/abc", nil, 1800000000)

	wantSig := "sha256=f53cf714f69187170c4fdb22c53e0b53578dcbcb61e63b56948d5e1fd8294a3e"
	wantTs := "1800000000"

	if sig != wantSig {
		t.Errorf("sig header = %q, want %q", sig, wantSig)
	}
	if ts != wantTs {
		t.Errorf("ts header = %q, want %q", ts, wantTs)
	}
}

func TestSignRequest_MethodCaseInsensitive(t *testing.T) {
	secret := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	sigLower, _ := SignRequest(secret, "post", "/api/v1/upload-sessions", []byte(`{"a":1}`), 1800000000)
	sigUpper, _ := SignRequest(secret, "POST", "/api/v1/upload-sessions", []byte(`{"a":1}`), 1800000000)

	if sigLower != sigUpper {
		t.Errorf("lowercase method sig %q != uppercase method sig %q", sigLower, sigUpper)
	}
}

// TestHTTPClient_SignsJSONRequest_WhenSecretConfigured drives an actual
// request through httpClient.request against an httptest server and confirms
// the signature headers are present and correctly formed when a signing
// secret is configured.
func TestHTTPClient_SignsJSONRequest_WhenSecretConfigured(t *testing.T) {
	const secret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	var gotSig, gotTs, gotPath, gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Aptr-Signature")
		gotTs = r.Header.Get("X-Aptr-Timestamp")
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	h := newHTTPClient(srv.URL, "aptr_test_token", secret)
	if err := h.request(context.Background(), http.MethodPost, "/api/v1/upload-sessions", map[string]int{"a": 1}, nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if gotSig == "" {
		t.Error("expected X-Aptr-Signature header to be set, got empty")
	}
	if gotTs == "" {
		t.Error("expected X-Aptr-Timestamp header to be set, got empty")
	}
	if gotPath != "/api/v1/upload-sessions" || gotMethod != http.MethodPost {
		t.Fatalf("unexpected request observed: method=%s path=%s", gotMethod, gotPath)
	}

	// The server received the exact body bytes `{"a":1}` (json.Marshal of
	// map[string]int{"a":1}), so the signature must match the known-answer
	// vector for that body.
	wantSig, wantTs := SignRequest(secret, "POST", "/api/v1/upload-sessions", []byte(`{"a":1}`), mustParseUnix(t, gotTs))
	if gotSig != wantSig {
		t.Errorf("sig header = %q, want %q", gotSig, wantSig)
	}
	if gotTs != wantTs {
		t.Errorf("ts header = %q, want %q", gotTs, wantTs)
	}
}

// TestHTTPClient_NoSignature_WhenSecretUnset confirms backwards compatibility:
// requests carry no signing headers when no signing secret is configured.
func TestHTTPClient_NoSignature_WhenSecretUnset(t *testing.T) {
	var gotSig, gotTs string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Aptr-Signature")
		gotTs = r.Header.Get("X-Aptr-Timestamp")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	h := newHTTPClient(srv.URL, "aptr_test_token", "")
	if err := h.request(context.Background(), http.MethodGet, "/api/v1/uploads/abc", nil, nil); err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if gotSig != "" {
		t.Errorf("expected no X-Aptr-Signature header, got %q", gotSig)
	}
	if gotTs != "" {
		t.Errorf("expected no X-Aptr-Timestamp header, got %q", gotTs)
	}
}

// TestHTTPClient_RequestRaw_Signs confirms requestRaw (used for binary
// downloads, always bodyless) is also signed when a secret is configured.
func TestHTTPClient_RequestRaw_Signs(t *testing.T) {
	const secret = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

	var gotSig, gotTs string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Aptr-Signature")
		gotTs = r.Header.Get("X-Aptr-Timestamp")
		w.Write([]byte("binary-data"))
	}))
	defer srv.Close()

	h := newHTTPClient(srv.URL, "aptr_test_token", secret)
	if _, err := h.requestRaw(context.Background(), http.MethodGet, "/api/v1/uploads/abc/qr"); err != nil {
		t.Fatalf("requestRaw failed: %v", err)
	}

	if gotSig == "" {
		t.Error("expected X-Aptr-Signature header to be set, got empty")
	}
	if gotTs == "" {
		t.Error("expected X-Aptr-Timestamp header to be set, got empty")
	}
}

func mustParseUnix(t *testing.T, s string) int64 {
	t.Helper()
	var v int64
	for _, c := range s {
		if c < '0' || c > '9' {
			t.Fatalf("invalid unix timestamp string: %q", s)
		}
		v = v*10 + int64(c-'0')
	}
	return v
}
