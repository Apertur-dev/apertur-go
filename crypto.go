package apertur

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// EncryptImage encrypts image data using AES-256-GCM with a random key, then
// wraps the AES key using RSA-OAEP with the provided PEM-encoded public key.
// The returned EncryptedPayload contains all fields base64-encoded, ready to be
// sent to the Apertur API.
func EncryptImage(imageData []byte, publicKeyPEM string) (*EncryptedPayload, error) {
	// Generate random AES-256 key (32 bytes) and GCM IV (12 bytes)
	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("apertur: failed to generate AES key: %w", err)
	}

	iv := make([]byte, 12)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("apertur: failed to generate IV: %w", err)
	}

	// Encrypt image data with AES-256-GCM
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to create GCM: %w", err)
	}

	// Seal appends the auth tag to the ciphertext
	encryptedData := gcm.Seal(nil, iv, imageData, nil)

	// Parse the RSA public key
	pemBlock, _ := pem.Decode([]byte(publicKeyPEM))
	if pemBlock == nil {
		return nil, fmt.Errorf("apertur: failed to decode PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to parse public key: %w", err)
	}

	rsaPub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("apertur: public key is not RSA")
	}

	// Wrap the AES key with RSA-OAEP (SHA-256)
	wrappedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPub, aesKey, nil)
	if err != nil {
		return nil, fmt.Errorf("apertur: failed to wrap AES key: %w", err)
	}

	return &EncryptedPayload{
		EncryptedKey:  base64.StdEncoding.EncodeToString(wrappedKey),
		IV:            base64.StdEncoding.EncodeToString(iv),
		EncryptedData: base64.StdEncoding.EncodeToString(encryptedData),
		Algorithm:     "RSA-OAEP+AES-256-GCM",
	}, nil
}
