package batch

import (
	"errors"
	"fmt"
)

const MaxIDs = 100

var ErrTooManyIDs = errors.New("too many ids")

func ParseIDs[ID comparable](raw []string, parse func(string) (ID, error)) ([]ID, error) {
	if len(raw) > MaxIDs {
		return nil, fmt.Errorf("%w: got %d, max %d", ErrTooManyIDs, len(raw), MaxIDs)
	}

	ids := make([]ID, 0, len(raw))
	seen := make(map[ID]bool, len(raw))
	for _, s := range raw {
		id, err := parse(s)
		if err != nil {
			return nil, err
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	return ids, nil
}
