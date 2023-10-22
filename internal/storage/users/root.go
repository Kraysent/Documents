package users

import (
	"context"
	"fmt"

	libstorage "documents/internal/library/storage"
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
	selector := libstorage.SqAnd(r.Fields)

	row, err := s.storage.QueryRowSq(ctx, sq.Select(ColumnID).
		From(TableName).
		Where(selector).
		PlaceholderFormat(sq.Dollar))
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
	row, err := s.storage.QueryRowSq(ctx, sq.Insert(TableName).
		Columns(ColumnUsername, ColumnGoogleID).
		Values(r.Username, r.GoogleID).
		Suffix(fmt.Sprintf("RETURNING %s", ColumnID)).
		PlaceholderFormat(sq.Dollar))
	if err != nil {
		return nil, err
	}

	var result CreateUserResult
	if err := row.Scan(&result.UserID); err != nil {
		return nil, err
	}

	return &result, nil
}
