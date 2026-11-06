package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func EncodeCursor[K any](key K) (string, error) {
	raw, err := json.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("encode cursor: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func DecodeCursor[K any](cursor string) (K, error) {
	var key K
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || len(raw) == 0 {
		return key, ErrInvalidCursor
	}
	if err := json.Unmarshal(raw, &key); err != nil {
		return key, ErrInvalidCursor
	}
	return key, nil
}
