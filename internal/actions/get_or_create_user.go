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
	var userID int64
	var status string

	res, err := repo.Storages.Users.GetUser(ctx, users.GetUserRequest{Fields: map[string]any{
		users.ColumnGoogleID: r.GoogleID,
	}})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, web.DatabaseError(err)
	}

	userID = res.UserID
	status = "found"

	if err != nil {
		res, err := repo.Storages.Users.CreateUser(ctx, users.CreateUserRequest{Username: r.Email, GoogleID: r.GoogleID})
		if err != nil {
			return nil, web.DatabaseError(err)
		}

		userID = res.UserID
		status = "created"
	}

	repo.SessionManager.Put(ctx, "user_id", userID)

	return &schema.GetOrCreateUserResponse{Status: status}, nil
}
