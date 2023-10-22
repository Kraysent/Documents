package actions

import (
	"context"
	"errors"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/users"
	"github.com/jackc/pgx/v5"
)

func GetOrCreateUser(
	ctx context.Context, repo *core.Repository, r schema.GetOrCreateUserRequest,
) (*schema.GetOrCreateUserResponse, error) {
	_, err := repo.Storages.Users.GetUser(ctx, users.GetUserRequest{Fields: map[string]any{
		users.ColumnGoogleID: r.GoogleID,
	}})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, web.DatabaseError(err)
	}
	if err != nil {
		_, err := repo.Storages.Users.CreateUser(ctx, users.CreateUserRequest{Username: r.Email, GoogleID: r.GoogleID})
		if err != nil {
			return nil, web.DatabaseError(err)
		}

		return &schema.GetOrCreateUserResponse{Status: "created"}, nil
	}

	return &schema.GetOrCreateUserResponse{Status: "found"}, nil
}
