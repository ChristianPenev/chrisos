package batch

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"testing"
)

var errBadID = errors.New("bad id")

func parseInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", errBadID, s)
	}
	return n, nil
}

func TestParseIDs(t *testing.T) {
	t.Run("parses and removes duplicates in order", func(t *testing.T) {
		got, err := ParseIDs([]string{"3", "1", "3", "2", "1"}, parseInt)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(got, []int{3, 1, 2}) {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		got, err := ParseIDs(nil, parseInt)
		if err != nil || len(got) != 0 {
			t.Fatalf("got %v, %v", got, err)
		}
	})

	t.Run("accepts exactly the maximum", func(t *testing.T) {
		raw := make([]string, MaxIDs)
		for i := range raw {
			raw[i] = strconv.Itoa(i)
		}
		got, err := ParseIDs(raw, parseInt)
		if err != nil || len(got) != MaxIDs {
			t.Fatalf("got %d ids, %v", len(got), err)
		}
	})

	t.Run("rejects more than the maximum", func(t *testing.T) {
		raw := make([]string, MaxIDs+1)
		for i := range raw {
			raw[i] = "1"
		}
		if _, err := ParseIDs(raw, parseInt); !errors.Is(err, ErrTooManyIDs) {
			t.Fatalf("err = %v, want %v", err, ErrTooManyIDs)
		}
	})

	t.Run("returns the parse error", func(t *testing.T) {
		if _, err := ParseIDs([]string{"1", "x"}, parseInt); !errors.Is(err, errBadID) {
			t.Fatalf("err = %v, want %v", err, errBadID)
		}
	})
}
