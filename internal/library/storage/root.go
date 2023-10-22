package storage

import (
	sq "github.com/Masterminds/squirrel"
)

func SqAnd(fields map[string]any) sq.And {
	selector := sq.And{}

	for key, value := range fields {
		selector = append(selector, sq.Eq{key: value})
	}

	return selector
}
