package apertur

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
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

// SignRequest computes the X-Aptr-Signature / X-Aptr-Timestamp header values
// for an outgoing API request.
// Calculation: HMAC-SHA256(`${timestamp}.${method}.${path}.${sha256hex(body)}`, secret).
//   - method is uppercased.
//   - path must be the exact request path sent to the transport, verbatim
//     (for apertur this includes the "/api/v1" prefix).
//   - body must be the exact bytes sent as the request body; nil hashes as
//     the empty string.
//
// It returns the ready-to-use header values: sigHeader is formatted as
// "sha256=<hex>" and tsHeader is the decimal unix timestamp.
func SignRequest(secret, method, path string, body []byte, timestamp int64) (sigHeader, tsHeader string) {
	bodyHash := sha256.Sum256(body)
	signatureBase := fmt.Sprintf("%d.%s.%s.%s", timestamp, strings.ToUpper(method), path, hex.EncodeToString(bodyHash[:]))

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signatureBase))
	signature := hex.EncodeToString(mac.Sum(nil))

	return "sha256=" + signature, strconv.FormatInt(timestamp, 10)
}
