package actions

import (
	"context"

	"documents/internal/actions/schema"
	"documents/internal/core"
	"documents/internal/library/web"
	"documents/internal/storage/users"
)

func CreateUser(ctx context.Context, repo *core.Repository, r schema.CreateUserRequest) (*schema.CreateUserResponse, error) {
	_, err := repo.Storages.Users.CreateUser(ctx, users.CreateUserRequest{Username: r.Username})
	if err != nil {
		return nil, web.DatabaseError(err)
	}

	return &schema.CreateUserResponse{Status: "success"}, nil
}
