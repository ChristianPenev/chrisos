package pagination

import (
	"errors"
	"fmt"
)

const (
	DefaultSize = 20
	MaxSize     = 100
)

var (
	ErrInvalidPageSize = errors.New("invalid page size")
	ErrInvalidCursor   = errors.New("invalid cursor")
)

type Request struct {
	Size  int
	After string
}

func NewRequest(size int, after string) (Request, error) {
	switch {
	case size < 0:
		return Request{}, fmt.Errorf("%w: must not be negative", ErrInvalidPageSize)
	case size == 0:
		size = DefaultSize
	case size > MaxSize:
		size = MaxSize
	}
	return Request{Size: size, After: after}, nil
}

type Page[T any] struct {
	Items      []T
	NextCursor string
}

func (p Page[T]) HasNextPage() bool {
	return p.NextCursor != ""
}
