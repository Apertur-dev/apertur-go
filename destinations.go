package apertur

import (
	"context"
	"fmt"
)

// DestinationsResource provides CRUD operations on delivery destinations.
type DestinationsResource struct {
	http *httpClient
}

// List returns all destinations for the given project.
func (d *DestinationsResource) List(ctx context.Context, projectID string) ([]Destination, error) {
	var result []Destination
	if err := d.http.request(ctx, "GET", fmt.Sprintf("/api/v1/projects/%s/destinations", projectID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create adds a new destination to the project.
func (d *DestinationsResource) Create(ctx context.Context, projectID string, config CreateDestinationConfig) (*Destination, error) {
	var result Destination
	if err := d.http.request(ctx, "POST", fmt.Sprintf("/api/v1/projects/%s/destinations", projectID), config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update modifies an existing destination.
func (d *DestinationsResource) Update(ctx context.Context, projectID, destID string, config UpdateDestinationConfig) (*Destination, error) {
	var result Destination
	path := fmt.Sprintf("/api/v1/projects/%s/destinations/%s", projectID, destID)
	if err := d.http.request(ctx, "PATCH", path, config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes a destination from the project.
func (d *DestinationsResource) Delete(ctx context.Context, projectID, destID string) error {
	path := fmt.Sprintf("/api/v1/projects/%s/destinations/%s", projectID, destID)
	return d.http.request(ctx, "DELETE", path, nil, nil)
}

// Test sends a test payload to the destination and returns the result.
func (d *DestinationsResource) Test(ctx context.Context, projectID, destID string) (*TestDestinationResult, error) {
	var result TestDestinationResult
	path := fmt.Sprintf("/api/v1/projects/%s/destinations/%s/test", projectID, destID)
	if err := d.http.request(ctx, "POST", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
