package apertur

// Config holds the configuration for the Apertur client.
type Config struct {
	// APIKey is the API key used for authentication (prefixed with "aptr_" or "aptr_test_").
	APIKey string
	// OAuthToken is an OAuth bearer token, used as an alternative to APIKey.
	OAuthToken string
	// BaseURL overrides the default API base URL. When empty, the URL is auto-detected
	// from the API key prefix: "aptr_test_" uses the sandbox, otherwise the live endpoint.
	BaseURL string
	// SigningSecret, when set, enables HMAC request signing: every JSON/string
	// request automatically carries X-Aptr-Signature / X-Aptr-Timestamp headers.
	// Leave empty to disable signing (default, backwards-compatible). Multipart
	// uploads are never signed.
	SigningSecret string
}

// --- Sessions ---

// CreateSessionOptions contains options for creating a new upload session.
type CreateSessionOptions struct {
	DestinationIDs    []string `json:"destination_ids,omitempty"`
	LongPolling       *bool    `json:"long_polling,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	ExpiresInHours    *int     `json:"expires_in_hours,omitempty"`
	ExpiresAt         string   `json:"expires_at,omitempty"`
	MaxImages         *int     `json:"max_images,omitempty"`
	AllowedMimeTypes  []string `json:"allowed_mime_types,omitempty"`
	MaxImageDimension *int     `json:"max_image_dimension,omitempty"`
	Password          string   `json:"password,omitempty"`
}

// QrSpecs describes the QR code endpoint and available parameters.
type QrSpecs struct {
	Endpoint string            `json:"endpoint"`
	Formats  []string          `json:"formats"`
	Params   map[string]string `json:"params"`
}

// SessionDestination describes a destination attached to a session.
type SessionDestination struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// CreateSessionResult is the response from creating an upload session.
type CreateSessionResult struct {
	UUID              string               `json:"uuid"`
	UploadURL         string               `json:"upload_url"`
	QRURL             string               `json:"qr_url"`
	ShortURL          string               `json:"short_url"`
	QRSpecs           QrSpecs              `json:"qr_specs"`
	Destinations      []SessionDestination `json:"destinations"`
	LongPolling       bool                 `json:"long_polling"`
	ExpiresAt         string               `json:"expires_at"`
	PasswordProtected bool                 `json:"password_protected"`
	Env               string               `json:"env"`
}

// Session represents the state of an upload session.
type Session struct {
	ID                        string   `json:"id"`
	Status                    string   `json:"status"`
	ExpiresAt                 string   `json:"expiresAt"`
	Tags                      []string `json:"tags"`
	ImagesPerSession          *int     `json:"imagesPerSession"`
	EffectiveMaxImages        *int     `json:"effectiveMaxImages"`
	EffectiveAllowedMimeTypes []string `json:"effectiveAllowedMimeTypes"`
	EffectiveMaxImageDimension *int    `json:"effectiveMaxImageDimension"`
	PasswordProtected         *bool    `json:"password_protected"`
	ServerPublicKey           string   `json:"serverPublicKey,omitempty"`
	E2EEnabled                *bool    `json:"e2eEnabled,omitempty"`
	E2EPublicKey              *string  `json:"e2ePublicKey,omitempty"`
	E2EDowngraded             *bool    `json:"e2eDowngraded,omitempty"`
}

// QrOptions configures the QR code rendering.
type QrOptions struct {
	Format      string `json:"format,omitempty"`
	Size        int    `json:"size,omitempty"`
	Style       string `json:"style,omitempty"`
	FG          string `json:"fg,omitempty"`
	BG          string `json:"bg,omitempty"`
	BorderSize  int    `json:"border_size,omitempty"`
	BorderColor string `json:"border_color,omitempty"`
}

// UpdateSessionOptions contains the fields that can be updated on a session.
type UpdateSessionOptions struct {
	ExpiresAt         string   `json:"expires_at,omitempty"`
	MaxImages         *int     `json:"max_images,omitempty"`
	AllowedMimeTypes  []string `json:"allowed_mime_types,omitempty"`
	MaxImageDimension *int     `json:"max_image_dimension,omitempty"`
	MaxImageSizeMB    *int     `json:"max_image_size_mb,omitempty"`
	Password          *string  `json:"password,omitempty"`
}

// SessionRow is a summary row returned in session list endpoints.
type SessionRow struct {
	ID                string   `json:"id"`
	CreatedAt         string   `json:"createdAt"`
	ExpiresAt         string   `json:"expiresAt"`
	Status            string   `json:"status"`
	ProjectID         string   `json:"projectId"`
	ProjectName       string   `json:"projectName"`
	ImagesCount       int      `json:"imagesCount"`
	ImagesDelivered   int      `json:"imagesDelivered"`
	ImagesFailed      int      `json:"imagesFailed"`
	DestinationsCount int      `json:"destinationsCount"`
	Tags              []string `json:"tags"`
	LongPollingEnabled bool    `json:"longPollingEnabled"`
	Label             *string  `json:"label"`
	Env               string   `json:"env"`
}

// SessionsListPage is a paginated response of session rows.
type SessionsListPage struct {
	Data       []SessionRow `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
	TotalPages int          `json:"totalPages"`
}

// ListParams holds common pagination parameters.
type ListParams struct {
	Page     *int `json:"page,omitempty"`
	PageSize *int `json:"pageSize,omitempty"`
}

// LimitParams holds a limit parameter for "recent" endpoints.
type LimitParams struct {
	Limit *int `json:"limit,omitempty"`
}

// VerifyPasswordResult is the response from verifying a session password.
type VerifyPasswordResult struct {
	Valid bool `json:"valid"`
}

// ShareSessionOptions contains options for sharing an upload session with a
// guest via email or SMS.
type ShareSessionOptions struct {
	// Channel is "email" or "sms".
	Channel string `json:"channel"`
	// Recipient is the destination email address or phone number.
	Recipient string `json:"recipient"`
	// SmsConsent must be true when Channel is "sms" (enforced server-side).
	SmsConsent *bool `json:"sms_consent,omitempty"`
	// Note is only used for email shares.
	Note string `json:"note,omitempty"`
}

// ShareResult is the response from sharing an upload session.
type ShareResult struct {
	Ok        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
	ShortURL  string `json:"short_url"`
}

// --- Upload ---

// UploadResult is the response from uploading an image.
type UploadResult struct {
	ID           string `json:"id"`
	Filename     string `json:"filename"`
	SizeBytes    int64  `json:"size_bytes"`
	Destinations int    `json:"destinations"`
	LongPolling  bool   `json:"long_polling"`
}

// UploadOptions contains optional parameters for image uploads.
type UploadOptions struct {
	Filename string
	MimeType string
	Source   string
	Password string
}

// --- Polling ---

// PollImage describes an image available for polling.
type PollImage struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	SizeBytes int64  `json:"size_bytes"`
	MimeType  string `json:"mime_type"`
	Source    string `json:"source"`
	CreatedAt string `json:"created_at"`
}

// PollResult is the response from the polling endpoint.
type PollResult struct {
	Images []PollImage `json:"images"`
}

// AckResult is the response from acknowledging a polled image.
type AckResult struct {
	Status string `json:"status"`
}

// PollProcessOptions configures the PollAndProcess loop.
type PollProcessOptions struct {
	// Interval is the delay between poll cycles. Defaults to 3 seconds.
	Interval int
}

// --- Delivery Status ---

// DeliveryDestinationStatus describes the delivery state for a single destination.
type DeliveryDestinationStatus struct {
	DestinationID string  `json:"destination_id"`
	Type          string  `json:"type"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	Attempts      int     `json:"attempts"`
	LastError     *string `json:"last_error"`
}

// DeliveryRecordStatus describes the delivery state for a single uploaded image.
type DeliveryRecordStatus struct {
	RecordID     string                      `json:"record_id"`
	Filename     string                      `json:"filename"`
	SizeBytes    int64                       `json:"size_bytes"`
	HasThumbnail bool                        `json:"has_thumbnail,omitempty"`
	Destinations []DeliveryDestinationStatus `json:"destinations"`
}

// DeliveryStatusResponse is the response from the delivery-status endpoint.
// The API returns the overall session status, the per-file delivery states,
// and the timestamp of the most recent change (used as the `pollFrom` cursor
// for long-polling follow-up requests).
type DeliveryStatusResponse struct {
	Status      string                 `json:"status"`
	Files       []DeliveryRecordStatus `json:"files"`
	LastChanged string                 `json:"lastChanged"`
}

// DeliveryStatusOptions contains optional parameters for the DeliveryStatus call.
type DeliveryStatusOptions struct {
	// PollFrom is an ISO 8601 timestamp. When set, the server will hold the
	// response for up to 5 minutes until the delivery state changes past this
	// cursor (new file, delivery transition, or session status change). Pass
	// the `LastChanged` value from a previous response to continue polling.
	//
	// Long-polling callers should pair this with a `context.Context` that has
	// at least a 6-minute deadline so the server releases the response first
	// under the happy path.
	PollFrom string
}

// --- Destinations ---

// Destination represents a configured delivery destination.
type Destination struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Config    map[string]interface{} `json:"config"`
	IsActive  bool                   `json:"isActive"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
}

// CreateDestinationConfig is the request body for creating a destination.
type CreateDestinationConfig struct {
	Type   string                 `json:"type"`
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// UpdateDestinationConfig is the request body for updating a destination.
type UpdateDestinationConfig struct {
	Name     string                 `json:"name,omitempty"`
	Config   map[string]interface{} `json:"config,omitempty"`
	IsActive *bool                  `json:"isActive,omitempty"`
}

// TestDestinationResult is the response from testing a destination.
type TestDestinationResult struct {
	Success bool    `json:"success"`
	Status  *int    `json:"status,omitempty"`
	Error   *string `json:"error,omitempty"`
	Message *string `json:"message,omitempty"`
}

// --- API Keys ---

// APIKey represents an API key resource.
type APIKey struct {
	ID                    string   `json:"id"`
	Prefix                string   `json:"prefix"`
	Label                 string   `json:"label"`
	Env                   string   `json:"env"`
	IsActive              bool     `json:"isActive"`
	LastUsedAt            *string  `json:"lastUsedAt"`
	MaxImages             *int     `json:"maxImages"`
	AllowedMimeTypes      []string `json:"allowedMimeTypes"`
	MaxImageDimension     *int     `json:"maxImageDimension"`
	LongPollingEnabled    bool     `json:"longPollingEnabled"`
	DefaultDestinations   []string `json:"defaultDestinations"`
	AllowedIPs            []string `json:"allowedIps"`
	AllowedDomains        []string `json:"allowedDomains"`
	TOTPEnabled           bool     `json:"totpEnabled"`
	ClientCertEnabled     bool     `json:"clientCertEnabled"`
	ClientCertFingerprint *string  `json:"clientCertFingerprint"`
	CreatedAt             string   `json:"createdAt"`
}

// CreateAPIKeyOptions contains options for creating a new API key.
type CreateAPIKeyOptions struct {
	Label             string   `json:"label"`
	MaxImages         *int     `json:"maxImages,omitempty"`
	AllowedMimeTypes  []string `json:"allowedMimeTypes,omitempty"`
	MaxImageDimension *int     `json:"maxImageDimension,omitempty"`
}

// CreateAPIKeyResult is the response from creating an API key. The PlainTextKey
// is only returned once and must be stored by the caller.
type CreateAPIKeyResult struct {
	Key          APIKey `json:"key"`
	PlainTextKey string `json:"plainTextKey"`
}

// UpdateAPIKeyOptions contains fields that can be updated on an API key.
type UpdateAPIKeyOptions struct {
	Label             string   `json:"label,omitempty"`
	IsActive          *bool    `json:"isActive,omitempty"`
	MaxImages         *int     `json:"maxImages,omitempty"`
	AllowedMimeTypes  []string `json:"allowedMimeTypes,omitempty"`
	MaxImageDimension *int     `json:"maxImageDimension,omitempty"`
	AllowedIPs        []string `json:"allowedIps,omitempty"`
	AllowedDomains    []string `json:"allowedDomains,omitempty"`
}

// KeyDestinationEntry describes a destination linked to a key.
type KeyDestinationEntry struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// KeyDestinations is the response from setting destinations on a key.
type KeyDestinations struct {
	Destinations       []KeyDestinationEntry `json:"destinations"`
	LongPollingEnabled bool                  `json:"longPollingEnabled"`
}

// SetKeyDestinationsRequest is the request body for setting a key's destinations.
type SetKeyDestinationsRequest struct {
	DestinationIDs     []string `json:"destination_ids"`
	LongPollingEnabled *bool    `json:"long_polling_enabled,omitempty"`
}

// --- Event Webhooks ---

// EventWebhook represents an event webhook subscription.
type EventWebhook struct {
	ID                  string            `json:"id"`
	ProjectID           string            `json:"projectId"`
	URL                 string            `json:"url"`
	Secret              string            `json:"secret"`
	SignatureMethod     string            `json:"signatureMethod"`
	Topics              []string          `json:"topics"`
	IsActive            bool              `json:"isActive"`
	MaxRetries          int               `json:"maxRetries"`
	RetryIntervals      []int             `json:"retryIntervals"`
	DisableAfterFailures int              `json:"disableAfterFailures"`
	ConsecutiveFailures int               `json:"consecutiveFailures"`
	CustomHeaders       map[string]string `json:"customHeaders"`
	DisabledAt          *string           `json:"disabledAt"`
	CreatedAt           string            `json:"createdAt"`
	UpdatedAt           string            `json:"updatedAt"`
}

// CreateEventWebhookConfig is the request body for creating an event webhook.
type CreateEventWebhookConfig struct {
	URL                  string            `json:"url"`
	Topics               []string          `json:"topics"`
	SignatureMethod      string            `json:"signatureMethod,omitempty"`
	MaxRetries           *int              `json:"maxRetries,omitempty"`
	RetryIntervals       []int             `json:"retryIntervals,omitempty"`
	DisableAfterFailures *int              `json:"disableAfterFailures,omitempty"`
	CustomHeaders        map[string]string `json:"customHeaders,omitempty"`
}

// UpdateEventWebhookConfig is the request body for updating an event webhook.
type UpdateEventWebhookConfig struct {
	URL                  string            `json:"url,omitempty"`
	Topics               []string          `json:"topics,omitempty"`
	IsActive             *bool             `json:"isActive,omitempty"`
	MaxRetries           *int              `json:"maxRetries,omitempty"`
	RetryIntervals       []int             `json:"retryIntervals,omitempty"`
	DisableAfterFailures *int              `json:"disableAfterFailures,omitempty"`
	CustomHeaders        map[string]string `json:"customHeaders,omitempty"`
}

// WebhookDelivery represents a single webhook delivery attempt.
type WebhookDelivery struct {
	ID           string  `json:"id"`
	EventLogID   string  `json:"eventLogId"`
	Topic        string  `json:"topic"`
	Status       string  `json:"status"`
	Attempts     int     `json:"attempts"`
	ResponseCode *int    `json:"responseCode"`
	ResponseBody *string `json:"responseBody"`
	DurationMs   int     `json:"durationMs"`
	LastError    *string `json:"lastError"`
	NextRetryAt  *string `json:"nextRetryAt"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
}

// WebhookDeliveriesResult is a paginated list of webhook deliveries.
type WebhookDeliveriesResult struct {
	Deliveries []WebhookDelivery `json:"deliveries"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

// WebhookDeliveriesOptions holds pagination options for listing deliveries.
type WebhookDeliveriesOptions struct {
	Page  *int `json:"page,omitempty"`
	Limit *int `json:"limit,omitempty"`
}

// WebhookTestResult is the response from testing a webhook endpoint.
type WebhookTestResult struct {
	Message string `json:"message"`
}

// WebhookRetryResult is the response from retrying a webhook delivery.
type WebhookRetryResult struct {
	Message string `json:"message"`
}

// --- Encryption ---

// ServerKey is the response from the encryption server-key endpoint.
type ServerKey struct {
	PublicKey string `json:"publicKey"`
}

// EncryptedPayload holds the encrypted image data and wrapped key.
type EncryptedPayload struct {
	EncryptedKey  string `json:"encryptedKey"`
	IV            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
	Algorithm     string `json:"algorithm"`
}

// --- Uploads ---

// UploadDestinationBreakdown describes the count of deliveries per destination type.
type UploadDestinationBreakdown struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

// UploadRow is a summary row for an uploaded image.
type UploadRow struct {
	ID                      string                       `json:"id"`
	Filename                string                       `json:"filename"`
	SizeBytes               int64                        `json:"sizeBytes"`
	MimeType                string                       `json:"mimeType"`
	Source                  string                       `json:"source"`
	IsEncrypted             bool                         `json:"isEncrypted"`
	Env                     string                       `json:"env"`
	CreatedAt               string                       `json:"createdAt"`
	SessionID               string                       `json:"sessionId"`
	ProjectID               string                       `json:"projectId"`
	ProjectName             string                       `json:"projectName"`
	DestinationsTotal       int                          `json:"destinationsTotal"`
	DestinationsDelivered   int                          `json:"destinationsDelivered"`
	DestinationsFailed      int                          `json:"destinationsFailed"`
	DestinationsBreakdown   []UploadDestinationBreakdown `json:"destinationsBreakdown"`
	Status                  string                       `json:"status"`
}

// UploadsListPage is a paginated response of upload rows.
type UploadsListPage struct {
	Data       []UploadRow `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

// --- Stats ---

// StatsTopProject describes a project with its session count.
type StatsTopProject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Sessions int    `json:"sessions"`
}

// Stats holds dashboard-level statistics.
type Stats struct {
	SessionsThisMonth   int               `json:"sessionsThisMonth"`
	SessionsTotal       int               `json:"sessionsTotal"`
	ImagesUploaded      int               `json:"imagesUploaded"`
	ImagesDelivered     int               `json:"imagesDelivered"`
	DeliverySuccessRate float64           `json:"deliverySuccessRate"`
	TotalProjects       int               `json:"totalProjects"`
	ActiveKeys          int               `json:"activeKeys"`
	TopProjects         []StatsTopProject `json:"topProjects"`
}
