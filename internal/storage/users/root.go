package users

import (
	"context"
	"fmt"

	libstorage "documents/internal/library/storage"
	"documents/internal/storage"
	sq "github.com/Masterminds/squirrel"
)

// UserStorage is a thin wrapper around users table.
type UserStorage interface {
	GetUser(context.Context, GetUserRequest) (*GetUserResult, error)
	CreateUser(context.Context, CreateUserRequest) (*CreateUserResult, error)
}

type userStorageImpl struct {
	storage storage.Storage
}

func NewUserStorage(store storage.Storage) *userStorageImpl {
	return &userStorageImpl{storage: store}
}

func (s *userStorageImpl) GetUser(ctx context.Context, r GetUserRequest) (*GetUserResult, error) {
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

func (s *userStorageImpl) CreateUser(ctx context.Context, r CreateUserRequest) (*CreateUserResult, error) {
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
