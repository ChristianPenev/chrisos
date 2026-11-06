package postgres

import (
	"errors"
	"os"
	"slices"
	"testing"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
)

type keysetItem struct {
	ID    int64 `gorm:"primaryKey"`
	Group string
}

func (keysetItem) TableName() string { return "keyset_test_items" }

func itemID(i keysetItem) int64 { return i.ID }

func newKeysetTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })

	if err := db.Exec("DROP TABLE IF EXISTS keyset_test_items").Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Exec(`CREATE TABLE keyset_test_items (id bigint PRIMARY KEY, "group" text NOT NULL)`).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec("DROP TABLE IF EXISTS keyset_test_items") })

	items := make([]keysetItem, 0, 7)
	for _, id := range []int64{5, 1, 7, 3, 2, 6, 4} {
		items = append(items, keysetItem{ID: id, Group: "a"})
	}
	items = append(items, keysetItem{ID: 100, Group: "b"})
	if err := db.Create(&items).Error; err != nil {
		t.Fatal(err)
	}
	return db
}

func ids(items []keysetItem) []int64 {
	out := make([]int64, len(items))
	for i, it := range items {
		out[i] = it.ID
	}
	return out
}

func TestKeyset(t *testing.T) {
	db := newKeysetTestDB(t)
	groupA := func() *gorm.DB { return db.Where(`"group" = ?`, "a") }

	t.Run("walks all pages in key order", func(t *testing.T) {
		var got [][]int64
		after := ""
		for range 10 {
			page, err := Keyset(groupA(), pagination.Request{Size: 3, After: after}, "id", itemID)
			if err != nil {
				t.Fatal(err)
			}
			got = append(got, ids(page.Items))
			if !page.HasNextPage() {
				break
			}
			after = page.NextCursor
		}

		want := [][]int64{{1, 2, 3}, {4, 5, 6}, {7}}
		if !slices.EqualFunc(got, want, slices.Equal) {
			t.Fatalf("pages = %v, want %v", got, want)
		}
	})

	t.Run("exact fit has no next page", func(t *testing.T) {
		page, err := Keyset(groupA(), pagination.Request{Size: 7}, "id", itemID)
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 7 || page.HasNextPage() {
			t.Fatalf("got %d items, next=%v", len(page.Items), page.HasNextPage())
		}
	})

	t.Run("inserted rows do not shift later pages", func(t *testing.T) {
		first, err := Keyset(groupA(), pagination.Request{Size: 3}, "id", itemID)
		if err != nil {
			t.Fatal(err)
		}
		if err := db.Exec(`INSERT INTO keyset_test_items (id, "group") VALUES (0, 'a')`).Error; err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { db.Exec("DELETE FROM keyset_test_items WHERE id = 0") })

		second, err := Keyset(groupA(), pagination.Request{Size: 3, After: first.NextCursor}, "id", itemID)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(ids(second.Items), []int64{4, 5, 6}) {
			t.Fatalf("second page = %v, want [4 5 6]", ids(second.Items))
		}
	})

	t.Run("zero size uses the default", func(t *testing.T) {
		page, err := Keyset(groupA(), pagination.Request{}, "id", itemID)
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 7 {
			t.Fatalf("got %d items", len(page.Items))
		}
	})

	t.Run("empty result", func(t *testing.T) {
		page, err := Keyset(db.Where(`"group" = ?`, "none"), pagination.Request{Size: 3}, "id", itemID)
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Items) != 0 || page.HasNextPage() {
			t.Fatalf("got %+v", page)
		}
	})

	invalid := map[string]string{
		"garbage":          "not a cursor",
		"wrong value type": mustCursor(t, map[string]string{"id": "abc"}),
		"other column":     mustCursor(t, map[string]int64{"created_at": 3}),
		"extra column":     mustCursor(t, map[string]int64{"id": 3, "x": 1}),
	}
	for name, cursor := range invalid {
		t.Run("rejects "+name+" cursor", func(t *testing.T) {
			_, err := Keyset(groupA(), pagination.Request{Size: 3, After: cursor}, "id", itemID)
			if !errors.Is(err, pagination.ErrInvalidCursor) {
				t.Fatalf("err = %v, want %v", err, pagination.ErrInvalidCursor)
			}
		})
	}

	t.Run("rejects negative size", func(t *testing.T) {
		_, err := Keyset(groupA(), pagination.Request{Size: -1}, "id", itemID)
		if !errors.Is(err, pagination.ErrInvalidPageSize) {
			t.Fatalf("err = %v, want %v", err, pagination.ErrInvalidPageSize)
		}
	})
}

func mustCursor[K any](t *testing.T, key K) string {
	t.Helper()
	c, err := pagination.EncodeCursor(key)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
