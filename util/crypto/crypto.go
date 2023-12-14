package crypto

import (
	"crypto/sha256"
	"encoding/base64"
)

func SHA256URLEncode(key string) string {
	h := sha256.Sum256([]byte(key))
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(h[:])
}
