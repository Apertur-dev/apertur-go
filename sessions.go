package apertur

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// SessionsResource provides methods for managing upload sessions.
type SessionsResource struct {
	http *httpClient
}

// Create creates a new upload session with the given options.
func (s *SessionsResource) Create(ctx context.Context, opts CreateSessionOptions) (*CreateSessionResult, error) {
	var result CreateSessionResult
	if err := s.http.request(ctx, "POST", "/api/v1/upload-sessions", opts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves the state of an upload session by UUID.
func (s *SessionsResource) Get(ctx context.Context, uuid string) (*Session, error) {
	var result Session
	if err := s.http.request(ctx, "GET", fmt.Sprintf("/api/v1/upload/%s/session", uuid), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update modifies an existing upload session.
func (s *SessionsResource) Update(ctx context.Context, uuid string, opts UpdateSessionOptions) (*Session, error) {
	var result Session
	if err := s.http.request(ctx, "PATCH", fmt.Sprintf("/api/v1/upload-sessions/%s", uuid), opts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// List returns a paginated list of sessions.
func (s *SessionsResource) List(ctx context.Context, params ListParams) (*SessionsListPage, error) {
	path := "/api/v1/sessions"
	qs := url.Values{}
	if params.Page != nil {
		qs.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PageSize != nil {
		qs.Set("pageSize", strconv.Itoa(*params.PageSize))
	}
	if encoded := qs.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var result SessionsListPage
	if err := s.http.request(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Recent returns the most recently created sessions.
func (s *SessionsResource) Recent(ctx context.Context, params LimitParams) ([]SessionRow, error) {
	path := "/api/v1/sessions/recent"
	if params.Limit != nil {
		path += "?limit=" + strconv.Itoa(*params.Limit)
	}

	var result []SessionRow
	if err := s.http.request(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// QR returns the QR code image bytes for the given session. The response format
// depends on the options (PNG, SVG, or JPEG).
func (s *SessionsResource) QR(ctx context.Context, uuid string, opts QrOptions) ([]byte, error) {
	qs := url.Values{}
	if opts.Format != "" {
		qs.Set("format", opts.Format)
	}
	if opts.Size > 0 {
		qs.Set("size", strconv.Itoa(opts.Size))
	}
	if opts.Style != "" {
		qs.Set("style", opts.Style)
	}
	if opts.FG != "" {
		qs.Set("fg", opts.FG)
	}
	if opts.BG != "" {
		qs.Set("bg", opts.BG)
	}
	if opts.BorderSize > 0 {
		qs.Set("borderSize", strconv.Itoa(opts.BorderSize))
	}
	if opts.BorderColor != "" {
		qs.Set("borderColor", opts.BorderColor)
	}

	path := fmt.Sprintf("/api/v1/upload-sessions/%s/qr", uuid)
	if encoded := qs.Encode(); encoded != "" {
		path += "?" + encoded
	}

	return s.http.requestRaw(ctx, "GET", path)
}

// VerifyPassword checks whether the given password is valid for the session.
func (s *SessionsResource) VerifyPassword(ctx context.Context, uuid, password string) (*VerifyPasswordResult, error) {
	body := map[string]string{"password": password}
	var result VerifyPasswordResult
	if err := s.http.request(ctx, "POST", fmt.Sprintf("/api/v1/upload/%s/verify-password", uuid), body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeliveryStatus returns the delivery status for all images in the session.
func (s *SessionsResource) DeliveryStatus(ctx context.Context, uuid string) ([]DeliveryRecordStatus, error) {
	var result []DeliveryRecordStatus
	if err := s.http.request(ctx, "GET", fmt.Sprintf("/api/v1/upload-sessions/%s/delivery-status", uuid), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}
