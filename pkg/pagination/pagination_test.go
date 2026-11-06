package pagination

import (
	"errors"
	"testing"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		wantSize int
		wantErr  error
	}{
		{name: "zero uses default", size: 0, wantSize: DefaultSize},
		{name: "keeps valid size", size: 50, wantSize: 50},
		{name: "keeps max size", size: MaxSize, wantSize: MaxSize},
		{name: "caps above max", size: MaxSize + 1, wantSize: MaxSize},
		{name: "rejects negative", size: -1, wantErr: ErrInvalidPageSize},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest(tt.size, "abc")
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("err = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Size != tt.wantSize || req.After != "abc" {
				t.Errorf("got %+v, want size %d after %q", req, tt.wantSize, "abc")
			}
		})
	}
}

func TestPageHasNextPage(t *testing.T) {
	if (Page[int]{Items: []int{1}}).HasNextPage() {
		t.Error("page without cursor reports a next page")
	}
	if !(Page[int]{Items: []int{1}, NextCursor: "x"}).HasNextPage() {
		t.Error("page with cursor reports no next page")
	}
}

type testKey struct {
	ID string `json:"id"`
}

func TestCursorRoundTrip(t *testing.T) {
	cursor, err := EncodeCursor(testKey{ID: "01a1b30a-a466-7df5-a40e-3f688a817abe"})
	if err != nil {
		t.Fatal(err)
	}

	got, err := DecodeCursor[testKey](cursor)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "01a1b30a-a466-7df5-a40e-3f688a817abe" {
		t.Errorf("got %+v", got)
	}
}

func TestCursorIsURLSafe(t *testing.T) {
	cursor, err := EncodeCursor(testKey{ID: "??>>??>>"})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range cursor {
		if c == '+' || c == '/' || c == '=' {
			t.Fatalf("cursor %q is not URL safe", cursor)
		}
	}
}

func TestDecodeCursorRejectsGarbage(t *testing.T) {
	for _, cursor := range []string{"", "not base64!", "bm90IGpzb24", "W10"} {
		if _, err := DecodeCursor[testKey](cursor); !errors.Is(err, ErrInvalidCursor) {
			t.Errorf("DecodeCursor(%q) err = %v, want %v", cursor, err, ErrInvalidCursor)
		}
	}
}
