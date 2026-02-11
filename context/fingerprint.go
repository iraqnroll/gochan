package context

import (
	"context"
)

type fKey string

const (
	fingerprintKey fKey = "fingerprint"
)

func StoreFingerprint(ctx context.Context, fingerprint string) context.Context {
	return context.WithValue(ctx, fingerprintKey, fingerprint)
}

func GetFingerprint(ctx context.Context) string {
	val := ctx.Value(fingerprintKey)
	fingerprint, ok := val.(string)

	if !ok {
		return ""
	}
	return fingerprint
}
