package rand

import (
	"crypto/rand"
	"encoding/base64"
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
