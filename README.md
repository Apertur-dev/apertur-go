# apertur-go

Official Go SDK for the [Apertur](https://apertur.ca) API. Collect photos from any mobile device via QR codes — no app required.

## Installation

```bash
go get github.com/Apertur-dev/apertur-go
```

Requires Go 1.21 or later. No external dependencies.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	apertur "github.com/Apertur-dev/apertur-go"
)

func main() {
	client := apertur.New(apertur.Config{APIKey: "aptr_xxxx"})
	ctx := context.Background()

	// Create a session and upload an image
	session, err := client.Sessions.Create(ctx, apertur.CreateSessionOptions{
		Tags: []string{"demo"},
	})
	if err != nil {
		panic(err)
	}

	f, _ := os.Open("./photo.jpg")
	defer f.Close()

	result, err := client.Upload.Image(ctx, session.UUID, f, apertur.UploadOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Uploaded: %s (%d bytes)\n", result.Filename, result.SizeBytes)
}
```

## Authentication

Two auth methods: API key (for server-to-server) and OAuth token (for third-party integrations).

```go
// API Key — environment auto-detected from prefix
client := apertur.New(apertur.Config{APIKey: "aptr_your_key_here"})

// Test/sandbox key
client := apertur.New(apertur.Config{APIKey: "aptr_test_your_key_here"})

// OAuth Token
client := apertur.New(apertur.Config{OAuthToken: "your_oauth_token"})

// Custom base URL (for self-hosted or staging)
client := apertur.New(apertur.Config{
	APIKey:  "aptr_xxx",
	BaseURL: "https://custom.api.example.com",
})
```

See [Authentication documentation](https://docs.apertur.ca/authentication)

## Sessions

Create upload sessions, check status, and manage password-protected sessions.

```go
ctx := context.Background()

// Create a session
longPolling := true
maxImages := 50
expiresIn := 48
session, err := client.Sessions.Create(ctx, apertur.CreateSessionOptions{
	Tags:           []string{"event-photos"},
	ExpiresInHours: &expiresIn,
	MaxImages:      &maxImages,
	Password:       "secret123",
	LongPolling:    &longPolling,
})

// Get session info
info, err := client.Sessions.Get(ctx, session.UUID)

// Update a session
newMax := 100
_, err = client.Sessions.Update(ctx, session.UUID, apertur.UpdateSessionOptions{
	MaxImages: &newMax,
})

// List sessions (paginated)
page := 1
list, err := client.Sessions.List(ctx, apertur.ListParams{Page: &page})

// Get recent sessions
limit := 5
recent, err := client.Sessions.Recent(ctx, apertur.LimitParams{Limit: &limit})

// Verify password for protected sessions
result, err := client.Sessions.VerifyPassword(ctx, session.UUID, "secret123")

// Get QR code as PNG bytes
qr, err := client.Sessions.QR(ctx, session.UUID, apertur.QrOptions{
	Format: "png",
	Size:   300,
})

// Check delivery status
status, err := client.Sessions.DeliveryStatus(ctx, session.UUID)
```

See [Sessions documentation](https://docs.apertur.ca/upload-sessions)

## Uploading Images

Upload images from any `io.Reader`. Supports optional server-side encryption.

```go
// Upload from file
f, _ := os.Open("./photo.jpg")
defer f.Close()
result, err := client.Upload.Image(ctx, sessionUUID, f, apertur.UploadOptions{
	Filename: "vacation.jpg",
	MimeType: "image/jpeg",
	Source:   "gallery",
})

// Upload with server-side encryption
key, _ := client.Encryption.GetServerKey(ctx)
f2, _ := os.Open("./photo.jpg")
defer f2.Close()
result, err = client.Upload.ImageEncrypted(ctx, sessionUUID, f2, key.PublicKey, apertur.UploadOptions{})
```

See [Upload documentation](https://docs.apertur.ca/upload-sessions)

## Long Polling

Retrieve uploaded images via polling instead of webhooks. Requires a session created with `LongPolling: true`.

```go
// Manual poll cycle
poll, err := client.Polling.List(ctx, sessionUUID)
for _, image := range poll.Images {
	data, err := client.Polling.Download(ctx, sessionUUID, image.ID)
	if err != nil { break }
	os.WriteFile("./downloads/"+image.Filename, data, 0644)
	client.Polling.Ack(ctx, sessionUUID, image.ID)
}

// Automatic poll + process loop (stops on context cancellation)
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

err = client.Polling.PollAndProcess(ctx, sessionUUID, func(image apertur.PollImage, data []byte) error {
	fmt.Printf("Received: %s (%d bytes)\n", image.Filename, image.SizeBytes)
	return os.WriteFile("./output/"+image.Filename, data, 0644)
}, apertur.PollProcessOptions{Interval: 3})

// Stop polling from another goroutine
cancel()
```

See [Long Polling documentation](https://docs.apertur.ca/long-polling)

## Receiving Webhooks

Verify webhook signatures to ensure payloads are authentic. Three verification methods match the Apertur signature schemes.

```go
import apertur "github.com/Apertur-dev/apertur-go"

// Image delivery webhook
isValid := apertur.VerifyWebhookSignature(bodyString, signatureHeader, webhookSecret)

// Event webhook (HMAC SHA256)
isValid = apertur.VerifyEventSignature(bodyString, timestampHeader, signatureHeader, eventSecret)

// Event webhook (Svix)
isValid = apertur.VerifySvixSignature(bodyString, svixID, svixTimestamp, svixSignature, eventSecret)
```

See [Webhook documentation](https://docs.apertur.ca/webhooks)

## Destinations

Manage delivery destinations (webhook, S3, Google Drive, etc.).

```go
destinations, err := client.Destinations.List(ctx, projectID)

webhook, err := client.Destinations.Create(ctx, projectID, apertur.CreateDestinationConfig{
	Type: "webhook",
	Name: "My Backend",
	Config: map[string]interface{}{
		"url":    "https://api.example.com/photos",
		"format": "json_base64",
	},
})

testResult, err := client.Destinations.Test(ctx, projectID, webhook.ID)
_, err = client.Destinations.Update(ctx, projectID, webhook.ID, apertur.UpdateDestinationConfig{
	Name: "Updated Name",
})
err = client.Destinations.Delete(ctx, projectID, webhook.ID)
```

See [Destinations documentation](https://docs.apertur.ca/destinations)

## API Keys

Manage API keys and their default destinations.

```go
keys, err := client.Keys.List(ctx, projectID)

created, err := client.Keys.Create(ctx, projectID, apertur.CreateAPIKeyOptions{
	Label: "Production",
})
fmt.Printf("Save this key: %s\n", created.PlainTextKey) // Only shown once!

// Set default destinations for a key
longPolling := true
keyDests, err := client.Keys.SetDestinations(ctx, created.Key.ID, []string{dest1ID, dest2ID}, &longPolling)
```

See [API Keys documentation](https://docs.apertur.ca/api-keys)

## Event Webhooks

Subscribe to project events (uploads, deliveries, billing changes, etc.).

```go
wh, err := client.Webhooks.Create(ctx, projectID, apertur.CreateEventWebhookConfig{
	URL:    "https://api.example.com/events",
	Topics: []string{"project.upload.*", "project.billing.plan_changed"},
})

// List deliveries
deliveries, err := client.Webhooks.Deliveries(ctx, projectID, wh.ID, apertur.WebhookDeliveriesOptions{})

// Retry a failed delivery
_, err = client.Webhooks.RetryDelivery(ctx, projectID, wh.ID, deliveries.Deliveries[0].ID)
```

See [Event Webhooks documentation](https://docs.apertur.ca/event-webhooks)

## Error Handling

All SDK errors implement the standard `error` interface. Use `errors.As()` for typed error handling.

```go
import "errors"

session, err := client.Sessions.Get(ctx, "nonexistent")
if err != nil {
	var rateLimitErr *apertur.RateLimitError
	var authErr *apertur.AuthenticationError
	var notFoundErr *apertur.NotFoundError
	var validationErr *apertur.ValidationError
	var apiErr *apertur.AperturError

	switch {
	case errors.As(err, &rateLimitErr):
		fmt.Printf("Rate limited. Retry after: %ds\n", rateLimitErr.RetryAfter)
	case errors.As(err, &authErr):
		fmt.Println("Invalid API key")
	case errors.As(err, &notFoundErr):
		fmt.Println("Resource not found")
	case errors.As(err, &validationErr):
		fmt.Printf("Validation error: %s\n", validationErr.Message)
	case errors.As(err, &apiErr):
		fmt.Printf("API error %d: %s (code: %s)\n", apiErr.StatusCode, apiErr.Message, apiErr.Code)
	default:
		fmt.Printf("Unexpected error: %v\n", err)
	}
}
```

## API Reference

For complete API documentation, visit [docs.apertur.ca](https://docs.apertur.ca).

## License

MIT
