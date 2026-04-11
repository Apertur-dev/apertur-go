package apertur

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// VerifyWebhookSignature verifies an image delivery webhook signature.
// The signature header is expected in the format "sha256=<hex>".
// Calculation: HMAC-SHA256(body, secret).
func VerifyWebhookSignature(body, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	expected := hex.EncodeToString(mac.Sum(nil))

	sig := signature
	if strings.HasPrefix(sig, "sha256=") {
		sig = sig[7:]
	}

	expectedBytes, err := hex.DecodeString(expected)
	if err != nil {
		return false
	}
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedBytes, sigBytes)
}

// VerifyEventSignature verifies an event webhook signature using HMAC SHA256.
// The signature header is expected in the format "sha256=<hex>".
// Calculation: HMAC-SHA256("${timestamp}.${body}", secret).
func VerifyEventSignature(body, timestamp, signature, secret string) bool {
	signatureBase := timestamp + "." + body

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signatureBase))
	expected := hex.EncodeToString(mac.Sum(nil))

	sig := signature
	if strings.HasPrefix(sig, "sha256=") {
		sig = sig[7:]
	}

	expectedBytes, err := hex.DecodeString(expected)
	if err != nil {
		return false
	}
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return false
	}

	return hmac.Equal(expectedBytes, sigBytes)
}

// VerifySvixSignature verifies an event webhook signature using the Svix method.
// The signature header is expected in the format "v1,<base64>".
// The secret must be hex-encoded.
// Calculation: HMAC-SHA256("${svixId}.${timestamp}.${body}", hex-decoded secret).
func VerifySvixSignature(body, svixID, timestamp, signature, secret string) bool {
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return false
	}

	signatureBase := svixID + "." + timestamp + "." + body

	mac := hmac.New(sha256.New, secretBytes)
	mac.Write([]byte(signatureBase))
	expected := mac.Sum(nil)

	sig := signature
	if strings.HasPrefix(sig, "v1,") {
		sig = sig[3:]
	}

	sigBytes, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}

	return hmac.Equal(expected, sigBytes)
}
