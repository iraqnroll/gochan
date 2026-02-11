package rand

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("rand.Bytes error : %w", err)
	}
	if nRead < n {
		return nil, fmt.Errorf("rand.Bytes error : Didn't read enough bytes")
	}

	return b, nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("rand.String error : %w", err)
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateFingerprint(ip, salt string) string {
	data := ip + salt
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}
