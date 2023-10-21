package users

import (
	"context"
	"fmt"

	"documents/internal/storage"
	sq "github.com/Masterminds/squirrel"
)

// UserStorage is a thin wrapper around users table.
type UserStorage struct {
	storage storage.Storage
}

func NewUserStorage(store storage.Storage) *UserStorage {
	return &UserStorage{storage: store}
}

func (s *UserStorage) GetUser(ctx context.Context, r GetUserRequest) (*GetUserResult, error) {
	q := sq.Select(ColumnID).
		From(TableName).
		Where(sq.Eq{ColumnUsername: r.Username}).
		PlaceholderFormat(sq.Dollar)

	row, err := s.storage.QueryRowSq(ctx, q)
	if err != nil {
		return nil, err
	}

	var result GetUserResult
	if err := row.Scan(&result.UserID); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *UserStorage) CreateUser(ctx context.Context, r CreateUserRequest) (*CreateUserResult, error) {
	q := sq.Insert(TableName).
		Columns(ColumnUsername).
		Values(r.Username).
		Suffix(fmt.Sprintf("RETURNING %s", ColumnID)).
		PlaceholderFormat(sq.Dollar)

	row, err := s.storage.QueryRowSq(ctx, q)
	if err != nil {
		return nil, err
	}

	var result CreateUserResult
	if err := row.Scan(&result.UserID); err != nil {
		return nil, err
	}

	return &result, nil
}
