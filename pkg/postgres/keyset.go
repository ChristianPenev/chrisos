package postgres

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ChristianPenev/chrisos/pkg/pagination"
)

func Keyset[T any, K any](db *gorm.DB, req pagination.Request, column string, key func(T) K) (pagination.Page[T], error) {
	req, err := pagination.NewRequest(req.Size, req.After)
	if err != nil {
		return pagination.Page[T]{}, err
	}

	col := clause.Column{Name: column}
	query := db
	if req.After != "" {
		after, err := decodeKeysetCursor[K](req.After, column)
		if err != nil {
			return pagination.Page[T]{}, err
		}
		query = query.Where("? > ?", col, after)
	}

	var rows []T
	err = query.Order(clause.OrderByColumn{Column: col}).Limit(req.Size + 1).Find(&rows).Error
	if err != nil {
		return pagination.Page[T]{}, err
	}

	if len(rows) <= req.Size {
		return pagination.Page[T]{Items: rows}, nil
	}

	rows = rows[:req.Size]
	next, err := pagination.EncodeCursor(map[string]K{column: key(rows[len(rows)-1])})
	if err != nil {
		return pagination.Page[T]{}, err
	}
	return pagination.Page[T]{Items: rows, NextCursor: next}, nil
}

func decodeKeysetCursor[K any](cursor, column string) (K, error) {
	var zero K
	values, err := pagination.DecodeCursor[map[string]K](cursor)
	if err != nil {
		return zero, err
	}
	value, ok := values[column]
	if !ok || len(values) != 1 {
		return zero, fmt.Errorf("%w: not a cursor for %q", pagination.ErrInvalidCursor, column)
	}
	return value, nil
}
