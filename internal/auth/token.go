package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

type Signer struct{ Secret []byte }

func (s Signer) Sign(sessionID, token string) string {
	mac := hmac.New(sha256.New, s.Secret)
	mac.Write([]byte(sessionID))
	mac.Write([]byte{0})
	mac.Write([]byte(token))
	sig := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString([]byte(sessionID + ":" + token + ":" + string(sig)))
}
func (s Signer) Verify(value string) (string, string, error) {
	raw, e := base64.RawURLEncoding.DecodeString(value)
	if e != nil {
		return "", "", e
	}
	parts := strings.SplitN(string(raw), ":", 3)
	if len(parts) != 3 {
		return "", "", errors.New("malformed token")
	}
	mac := hmac.New(sha256.New, s.Secret)
	mac.Write([]byte(parts[0]))
	mac.Write([]byte{0})
	mac.Write([]byte(parts[1]))
	if !hmac.Equal(mac.Sum(nil), []byte(parts[2])) {
		return "", "", errors.New("invalid signature")
	}
	return parts[0], parts[1], nil
}
