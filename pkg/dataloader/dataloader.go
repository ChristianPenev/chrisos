package dataloader

import (
	"context"
	"errors"
	"fmt"

	"github.com/vikstrous/dataloadgen"

	"github.com/ChristianPenev/chrisos/pkg/batch"
)

var ErrNotFound = errors.New("not found")

type Loader[K comparable, V any] = dataloadgen.Loader[K, V]

func New[K comparable, V any](fetch func(ctx context.Context, keys []K) ([]V, error), key func(V) K) *Loader[K, V] {
	return dataloadgen.NewLoader(func(ctx context.Context, keys []K) ([]V, []error) {
		items, err := fetch(ctx, keys)
		if err != nil {
			return nil, []error{err}
		}
		return byKey(keys, items, key)
	}, dataloadgen.WithBatchCapacity(batch.MaxIDs))
}

func byKey[K comparable, V any](keys []K, items []V, key func(V) K) ([]V, []error) {
	found := make(map[K]V, len(items))
	for _, item := range items {
		found[key(item)] = item
	}

	values := make([]V, len(keys))
	var errs []error
	for i, k := range keys {
		item, ok := found[k]
		if !ok {
			if errs == nil {
				errs = make([]error, len(keys))
			}
			errs[i] = fmt.Errorf("%w: %v", ErrNotFound, k)
			continue
		}
		values[i] = item
	}
	return values, errs
}
