package clob

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"github.com/darmayasa221/polymarket-go/internal/commons/timeutil"
)

// BuildL2Signature computes the L2 request signature:
//
//	HMAC-SHA256(base64-decoded(apiSecret), timestamp+method+requestPath+body)
//
// Returns the result as a base64-encoded string.
func BuildL2Signature(apiSecret, timestamp, method, requestPath, body string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(apiSecret)
	if err != nil {
		return "", fmt.Errorf("clob: decode api secret: %w", err)
	}
	message := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// setL2Headers adds the 5 required L2 auth headers to req.
func setL2Headers(req *http.Request, cfg Config, body string) error {
	timestamp := strconv.FormatInt(timeutil.Now().Unix(), 10)
	sig, err := BuildL2Signature(cfg.APISecret, timestamp, req.Method, req.URL.Path, body)
	if err != nil {
		return err
	}
	req.Header.Set("POLY_ADDRESS", cfg.Address)
	req.Header.Set("POLY_SIGNATURE", sig)
	req.Header.Set("POLY_TIMESTAMP", timestamp)
	req.Header.Set("POLY_API_KEY", cfg.APIKey)
	req.Header.Set("POLY_PASSPHRASE", cfg.APIPassphrase)
	return nil
}
