package dataloader

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"

	"github.com/ChristianPenev/chrisos/pkg/batch"
)

type item struct {
	ID   int
	Name string
}

type fakeSource struct {
	mu      sync.Mutex
	items   map[int]item
	batches [][]int
	err     error
}

func (s *fakeSource) fetch(_ context.Context, ids []int) ([]item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.batches = append(s.batches, slices.Clone(ids))
	if s.err != nil {
		return nil, s.err
	}

	var out []item
	for i := len(ids) - 1; i >= 0; i-- {
		if it, ok := s.items[ids[i]]; ok {
			out = append(out, it)
		}
	}
	return out, nil
}

func newSource(n int) *fakeSource {
	s := &fakeSource{items: map[int]item{}}
	for i := 1; i <= n; i++ {
		s.items[i] = item{ID: i}
	}
	return s
}

func itemID(i item) int { return i.ID }

func loadConcurrently(l *Loader[int, item], ids []int) ([]item, []error) {
	values := make([]item, len(ids))
	errs := make([]error, len(ids))
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Go(func() {
			values[i], errs[i] = l.Load(context.Background(), id)
		})
	}
	wg.Wait()
	return values, errs
}

func TestLoaderBatchesConcurrentLoads(t *testing.T) {
	src := newSource(3)
	values, errs := loadConcurrently(New(src.fetch, itemID), []int{1, 2, 1, 3})

	for i, err := range errs {
		if err != nil {
			t.Fatalf("load %d: %v", i, err)
		}
	}
	if got := []int{values[0].ID, values[1].ID, values[2].ID, values[3].ID}; !slices.Equal(got, []int{1, 2, 1, 3}) {
		t.Fatalf("values matched to wrong keys: %v", got)
	}
	if len(src.batches) != 1 {
		t.Fatalf("want 1 batch, got %d: %v", len(src.batches), src.batches)
	}
	if b := slices.Sorted(slices.Values(src.batches[0])); !slices.Equal(b, []int{1, 2, 3}) {
		t.Fatalf("batch keys = %v, want each key once", src.batches[0])
	}
}

func TestLoaderMissingKeysFailAlone(t *testing.T) {
	src := newSource(2)
	values, errs := loadConcurrently(New(src.fetch, itemID), []int{1, 99, 2})

	if errs[0] != nil || errs[2] != nil {
		t.Fatalf("found keys failed: %v", errs)
	}
	if values[0].ID != 1 || values[2].ID != 2 {
		t.Fatalf("got %+v", values)
	}
	if !errors.Is(errs[1], ErrNotFound) {
		t.Fatalf("missing key err = %v, want %v", errs[1], ErrNotFound)
	}
}

func TestLoaderSingleMissingKey(t *testing.T) {
	_, err := New(newSource(1).fetch, itemID).Load(context.Background(), 99)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want %v", err, ErrNotFound)
	}
}

func TestLoaderFetchErrorFailsTheBatch(t *testing.T) {
	boom := errors.New("user-service unavailable")
	src := newSource(2)
	src.err = boom

	_, errs := loadConcurrently(New(src.fetch, itemID), []int{1, 2})
	for i, err := range errs {
		if !errors.Is(err, boom) {
			t.Fatalf("load %d err = %v, want %v", i, err, boom)
		}
	}
}

func TestLoaderSplitsBatchesAtMaxIDs(t *testing.T) {
	src := newSource(batch.MaxIDs + 50)
	ids := make([]int, batch.MaxIDs+50)
	for i := range ids {
		ids[i] = i + 1
	}

	_, errs := loadConcurrently(New(src.fetch, itemID), ids)
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(src.batches) != 2 {
		t.Fatalf("want 2 batches, got %d", len(src.batches))
	}
	for _, b := range src.batches {
		if len(b) > batch.MaxIDs {
			t.Fatalf("batch of %d keys exceeds %d", len(b), batch.MaxIDs)
		}
	}
}

func TestLoaderCachesWithinItsLifetime(t *testing.T) {
	src := newSource(1)
	l := New(src.fetch, itemID)

	for range 3 {
		if _, err := l.Load(context.Background(), 1); err != nil {
			t.Fatal(err)
		}
	}
	if len(src.batches) != 1 {
		t.Fatalf("want 1 fetch, got %d", len(src.batches))
	}
}
