package apertur

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// UploadResource provides methods for uploading images to a session.
type UploadResource struct {
	http *httpClient
}

// Image uploads an image to the session identified by uuid. The file parameter
// is an io.Reader containing the image data. Use UploadOptions to set the
// filename, MIME type, source, and session password.
func (u *UploadResource) Image(ctx context.Context, uuid string, file io.Reader, opts UploadOptions) (*UploadResult, error) {
	filename := opts.Filename
	if filename == "" {
		filename = "image.jpg"
	}
	mimeType := opts.MimeType
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create the file part with the correct Content-Type
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, filename))
	partHeader.Set("Content-Type", mimeType)
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create multipart part: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("apertur: failed to write file data: %w", err)
	}

	if opts.Source != "" {
		if err := writer.WriteField("source", opts.Source); err != nil {
			return nil, fmt.Errorf("apertur: failed to write source field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("apertur: failed to close multipart writer: %w", err)
	}

	reqURL := u.http.url(fmt.Sprintf("/api/v1/upload/%s/images", uuid))
	req, err := http.NewRequest("POST", reqURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	if opts.Password != "" {
		req.Header.Set("x-session-password", opts.Password)
	}

	var result UploadResult
	if err := u.http.requestCustom(ctx, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// encryptedUploadPayload is the JSON body sent for encrypted uploads.
type encryptedUploadPayload struct {
	EncryptedKey  string `json:"encryptedKey"`
	IV            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
	Algorithm     string `json:"algorithm"`
	Filename      string `json:"filename"`
	MimeType      string `json:"mimeType"`
	Source        string `json:"source"`
}

// ImageEncrypted uploads an image encrypted with the server's public key.
// The image data is read from file, encrypted using EncryptImage, and sent as
// a JSON payload with the X-Aptr-Encrypted header.
func (u *UploadResource) ImageEncrypted(ctx context.Context, uuid string, file io.Reader, publicKeyPEM string, opts UploadOptions) (*UploadResult, error) {
	imageData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to read file data: %w", err)
	}

	encrypted, err := EncryptImage(imageData, publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to encrypt image: %w", err)
	}

	filename := opts.Filename
	if filename == "" {
		filename = "image.jpg"
	}
	mimeType := opts.MimeType
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	source := opts.Source
	if source == "" {
		source = "sdk"
	}

	payload := encryptedUploadPayload{
		EncryptedKey:  encrypted.EncryptedKey,
		IV:            encrypted.IV,
		EncryptedData: encrypted.EncryptedData,
		Algorithm:     encrypted.Algorithm,
		Filename:      filename,
		MimeType:      mimeType,
		Source:        source,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to marshal encrypted payload: %w", err)
	}

	reqURL := u.http.url(fmt.Sprintf("/api/v1/upload/%s/images", uuid))
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Aptr-Encrypted", "default")
	if opts.Password != "" {
		req.Header.Set("x-session-password", opts.Password)
	}

	var result UploadResult
	if err := u.http.requestCustom(ctx, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
