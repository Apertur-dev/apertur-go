package apertur

import "context"

// EncryptionResource provides access to the server's encryption keys.
type EncryptionResource struct {
	http *httpClient
}

// GetServerKey retrieves the server's public key used for end-to-end encryption.
func (e *EncryptionResource) GetServerKey(ctx context.Context) (*ServerKey, error) {
	var result ServerKey
	if err := e.http.request(ctx, "GET", "/api/v1/encryption/server-key", nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
