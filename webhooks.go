package apertur

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// WebhooksResource provides CRUD operations on event webhooks.
type WebhooksResource struct {
	http *httpClient
}

// List returns all event webhooks for the given project.
func (w *WebhooksResource) List(ctx context.Context, projectID string) ([]EventWebhook, error) {
	var result []EventWebhook
	if err := w.http.request(ctx, "GET", fmt.Sprintf("/api/v1/projects/%s/webhooks", projectID), nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Create registers a new event webhook for the project.
func (w *WebhooksResource) Create(ctx context.Context, projectID string, config CreateEventWebhookConfig) (*EventWebhook, error) {
	var result EventWebhook
	if err := w.http.request(ctx, "POST", fmt.Sprintf("/api/v1/projects/%s/webhooks", projectID), config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Update modifies an existing event webhook.
func (w *WebhooksResource) Update(ctx context.Context, projectID, webhookID string, config UpdateEventWebhookConfig) (*EventWebhook, error) {
	var result EventWebhook
	path := fmt.Sprintf("/api/v1/projects/%s/webhooks/%s", projectID, webhookID)
	if err := w.http.request(ctx, "PATCH", path, config, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete removes an event webhook from the project.
func (w *WebhooksResource) Delete(ctx context.Context, projectID, webhookID string) error {
	path := fmt.Sprintf("/api/v1/projects/%s/webhooks/%s", projectID, webhookID)
	return w.http.request(ctx, "DELETE", path, nil, nil)
}

// Test sends a test event to the webhook endpoint.
func (w *WebhooksResource) Test(ctx context.Context, projectID, webhookID string) (*WebhookTestResult, error) {
	var result WebhookTestResult
	path := fmt.Sprintf("/api/v1/projects/%s/webhooks/%s/test", projectID, webhookID)
	if err := w.http.request(ctx, "POST", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Deliveries returns a paginated list of delivery attempts for the webhook.
func (w *WebhooksResource) Deliveries(ctx context.Context, projectID, webhookID string, opts WebhookDeliveriesOptions) (*WebhookDeliveriesResult, error) {
	path := fmt.Sprintf("/api/v1/projects/%s/webhooks/%s/deliveries", projectID, webhookID)
	qs := url.Values{}
	if opts.Page != nil {
		qs.Set("page", strconv.Itoa(*opts.Page))
	}
	if opts.Limit != nil {
		qs.Set("limit", strconv.Itoa(*opts.Limit))
	}
	if encoded := qs.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var result WebhookDeliveriesResult
	if err := w.http.request(ctx, "GET", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RetryDelivery retries a failed webhook delivery.
func (w *WebhooksResource) RetryDelivery(ctx context.Context, projectID, webhookID, deliveryID string) (*WebhookRetryResult, error) {
	var result WebhookRetryResult
	path := fmt.Sprintf("/api/v1/projects/%s/webhooks/%s/deliveries/%s/retry", projectID, webhookID, deliveryID)
	if err := w.http.request(ctx, "POST", path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
