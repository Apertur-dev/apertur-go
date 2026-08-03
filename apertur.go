// Package apertur provides the official Go client for the Apertur API.
//
// Apertur lets you collect photos from any mobile device via QR codes — no app
// required. This SDK wraps the REST API with idiomatic Go types, context-aware
// methods, and built-in error handling.
//
// Create a client with New and access resources through the named fields:
//
//	client := apertur.New(apertur.Config{APIKey: "aptr_xxxx"})
//	session, err := client.Sessions.Create(ctx, apertur.CreateSessionOptions{})
package apertur

import (
	"errors"
	"strings"
)

const (
	defaultBaseURL = "https://api.aptr.ca"
	sandboxBaseURL = "https://sandbox.api.aptr.ca"
)

// Apertur is the top-level client for the Apertur API. Construct one with New.
// All resource operations are accessed through the exported fields.
type Apertur struct {
	// Env indicates the target environment, either "live" or "test".
	// It is inferred from the API key prefix unless BaseURL is overridden.
	Env string

	// Sessions provides operations on upload sessions.
	Sessions *SessionsResource
	// Upload provides image upload operations.
	Upload *UploadResource
	// Uploads provides listing operations for uploaded images.
	Uploads *UploadsResource
	// Polling provides long-polling operations to retrieve uploaded images.
	Polling *PollingResource
	// Destinations provides CRUD operations on delivery destinations.
	Destinations *DestinationsResource
	// Keys provides CRUD operations on API keys.
	Keys *KeysResource
	// Webhooks provides CRUD operations on event webhooks.
	Webhooks *WebhooksResource
	// Encryption provides access to server encryption keys.
	Encryption *EncryptionResource
	// Stats provides dashboard statistics.
	Stats *StatsResource

	http *httpClient
}

// New creates a new Apertur client with the given configuration. The Config must
// include either an APIKey or OAuthToken; otherwise New panics.
//
// The target environment is auto-detected from the key prefix: keys starting
// with "aptr_test_" target the sandbox, all others target the live API. Set
// Config.BaseURL to override this behaviour.
func New(cfg Config) *Apertur {
	token := cfg.APIKey
	if token == "" {
		token = cfg.OAuthToken
	}
	if token == "" {
		panic(errors.New("apertur: either APIKey or OAuthToken must be provided"))
	}

	env := "live"
	if strings.HasPrefix(token, "aptr_test_") {
		env = "test"
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		if env == "test" {
			baseURL = sandboxBaseURL
		} else {
			baseURL = defaultBaseURL
		}
	}

	h := newHTTPClient(baseURL, token, cfg.SigningSecret)

	a := &Apertur{
		Env:  env,
		http: h,
	}

	a.Sessions = &SessionsResource{http: h}
	a.Upload = &UploadResource{http: h}
	a.Uploads = &UploadsResource{http: h}
	a.Polling = &PollingResource{http: h}
	a.Destinations = &DestinationsResource{http: h}
	a.Keys = &KeysResource{http: h}
	a.Webhooks = &WebhooksResource{http: h}
	a.Encryption = &EncryptionResource{http: h}
	a.Stats = &StatsResource{http: h}

	return a
}
