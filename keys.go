package apertur

import (
	"context"
	"fmt"
)

// KeysResource provides CRUD operations on API keys.
type KeysResource struct {
	http *httpClient
}

// List returns all API keys for the given project.
func (k *KeysResource) List(ctx context.Context, projectID string) ([]APIKey, error) {
	var result []APIKey
	if err := k.http.request(ctx, "GET", fmt.Sprintf("/api/v1/projects/%s/keys", projectID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create generates a new API key for the project. The plain-text key is only
// available in the response and must be stored by the caller.
func (k *KeysResource) Create(ctx context.Context, projectID string, opts CreateAPIKeyOptions) (*CreateAPIKeyResult, error) {
	var result CreateAPIKeyResult
	if err := k.http.request(ctx, "POST", fmt.Sprintf("/api/v1/projects/%s/keys", projectID), opts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update modifies an existing API key.
func (k *KeysResource) Update(ctx context.Context, projectID, keyID string, opts UpdateAPIKeyOptions) (*APIKey, error) {
	var result APIKey
	path := fmt.Sprintf("/api/v1/projects/%s/keys/%s", projectID, keyID)
	if err := k.http.request(ctx, "PATCH", path, opts, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an API key from the project.
func (k *KeysResource) Delete(ctx context.Context, projectID, keyID string) error {
	path := fmt.Sprintf("/api/v1/projects/%s/keys/%s", projectID, keyID)
	return k.http.request(ctx, "DELETE", path, nil, nil)
}

// SetDestinations configures the default destinations and long-polling setting
// for the given key. Pass nil for longPolling to leave the setting unchanged.
func (k *KeysResource) SetDestinations(ctx context.Context, keyID string, destinationIDs []string, longPolling *bool) (*KeyDestinations, error) {
	body := SetKeyDestinationsRequest{
		DestinationIDs:     destinationIDs,
		LongPollingEnabled: longPolling,
	}
	var result KeyDestinations
	path := fmt.Sprintf("/api/v1/keys/%s/destinations", keyID)
	if err := k.http.request(ctx, "PUT", path, body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
