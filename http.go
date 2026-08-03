package apertur

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// httpClient wraps the standard net/http client with auth and JSON handling.
type httpClient struct {
	baseURL       string
	authHeader    string
	signingSecret string
	client        *http.Client
}

// newHTTPClient creates an httpClient with the given base URL, bearer token,
// and optional HMAC signing secret. An empty signingSecret disables request
// signing (backwards-compatible default).
func newHTTPClient(baseURL, token, signingSecret string) *httpClient {
	return &httpClient{
		baseURL:       strings.TrimRight(baseURL, "/"),
		authHeader:    "Bearer " + token,
		signingSecret: signingSecret,
		client:        &http.Client{},
	}
}

// applySignature sets the X-Aptr-Signature / X-Aptr-Timestamp headers on req
// when a signing secret is configured; it is a no-op otherwise. body must be
// the exact bytes that will be sent on the wire (nil for no body) — the
// server hashes whatever it actually receives, so an approximation here
// would just produce a signature that fails verification. This is only used
// for the JSON/string request path: multipart uploads (requestCustom) are
// never signed, since the multipart envelope (boundary, part framing) is
// serialized by the transport after this point and we never see the exact
// bytes it puts on the wire.
func (h *httpClient) applySignature(req *http.Request, method, path string, body []byte) {
	if h.signingSecret == "" {
		return
	}
	sig, ts := SignRequest(h.signingSecret, method, path, body, time.Now().Unix())
	req.Header.Set("X-Aptr-Signature", sig)
	req.Header.Set("X-Aptr-Timestamp", ts)
}

// request performs a JSON API request. If body is non-nil it is marshalled as
// JSON and sent. The response is unmarshalled into result. For 204 No Content,
// result is left untouched and nil is returned.
func (h *httpClient) request(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("apertur: failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, h.baseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("apertur: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", h.authHeader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	h.applySignature(req, method, path, data)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("apertur: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return h.handleError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("apertur: failed to decode response: %w", err)
		}
	}

	return nil
}

// requestRaw performs a request and returns the raw response bytes.
// This is used for binary endpoints such as QR codes and image downloads.
func (h *httpClient) requestRaw(ctx context.Context, method, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, h.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create request: %w", err)
	}

	req.Header.Set("Authorization", h.authHeader)
	h.applySignature(req, method, path, nil)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("apertur: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, h.handleError(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to read response body: %w", err)
	}

	return data, nil
}

// requestCustom performs a request with a pre-built *http.Request. The caller is
// responsible for setting the method, body, and content-type. Authorization is
// added automatically. The response is unmarshalled into result.
func (h *httpClient) requestCustom(ctx context.Context, req *http.Request, result interface{}) error {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", h.authHeader)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("apertur: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return h.handleError(resp)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("apertur: failed to decode response: %w", err)
		}
	}

	return nil
}

// handleError parses an error response and returns the appropriate typed error.
func (h *httpClient) handleError(resp *http.Response) error {
	var errBody struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
		errBody.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	if errBody.Message == "" {
		errBody.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	switch resp.StatusCode {
	case 401:
		return NewAuthenticationError(errBody.Message)
	case 404:
		return NewNotFoundError(errBody.Message)
	case 429:
		retryAfter := 0
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}
		return NewRateLimitError(errBody.Message, retryAfter)
	case 400:
		return NewValidationError(errBody.Message)
	default:
		return &AperturError{
			StatusCode: resp.StatusCode,
			Code:       errBody.Code,
			Message:    errBody.Message,
		}
	}
}

// url returns the full URL for the given path.
func (h *httpClient) url(path string) string {
	return h.baseURL + path
}
