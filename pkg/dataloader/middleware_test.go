package dataloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testLoaders struct{ id int }

func TestMiddlewareCreatesLoadersPerRequest(t *testing.T) {
	created := 0
	var seen []*testLoaders

	handler := Middleware(func() *testLoaders {
		created++
		return &testLoaders{id: created}
	})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, For[*testLoaders](r.Context()))
	}))

	for range 2 {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, "/graphql", nil))
	}

	if created != 2 || len(seen) != 2 {
		t.Fatalf("created %d loaders, handlers saw %d", created, len(seen))
	}
	if seen[0] == seen[1] {
		t.Fatal("two requests shared the same loaders")
	}
}

func TestForPanicsWithoutMiddleware(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("For did not panic")
		}
	}()
	For[*testLoaders](context.Background())
}
