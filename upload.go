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

// encryptedUploadEnvelope is the JSON envelope sent as the multipart file part
// bytes for default-mode encrypted uploads. The server reads the multipart file,
// parses these bytes as JSON, then decrypts with its private key. Only these four
// camelCase fields are expected — filename and source travel as multipart metadata.
type encryptedUploadEnvelope struct {
	EncryptedKey  string `json:"encryptedKey"`
	IV            string `json:"iv"`
	EncryptedData string `json:"encryptedData"`
	Algorithm     string `json:"algorithm"`
}

// ImageEncrypted uploads an image encrypted with the server's public key.
// The image data is read from file and encrypted using EncryptImage. The
// resulting envelope is sent as the bytes of a multipart "file" part with the
// X-Aptr-Encrypted: default header; the server decrypts it server-side.
func (u *UploadResource) ImageEncrypted(ctx context.Context, uuid string, file io.Reader, publicKeyPEM string, opts UploadOptions) (*UploadResult, error) {
	imageData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to read file data: %w", err)
	}

	encrypted, err := EncryptImage(imageData, publicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to encrypt image: %w", err)
	}

	envelope := encryptedUploadEnvelope{
		EncryptedKey:  encrypted.EncryptedKey,
		IV:            encrypted.IV,
		EncryptedData: encrypted.EncryptedData,
		Algorithm:     encrypted.Algorithm,
	}

	envelopeBytes, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to marshal encrypted envelope: %w", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// The encrypted envelope JSON is sent as the bytes of the multipart "file"
	// part. The server reads request.file(), parses these bytes as the encrypted
	// payload, then decrypts. The mime is octet-stream (the real image mime is
	// re-detected server-side after decryption).
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="file"; filename="image.jpg.enc"`)
	partHeader.Set("Content-Type", "application/octet-stream")
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create multipart part: %w", err)
	}

	if _, err := part.Write(envelopeBytes); err != nil {
		return nil, fmt.Errorf("apertur: failed to write envelope data: %w", err)
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
