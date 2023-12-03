package actions

import (
	"context"

	"documents/internal/core"
)

type Request interface {
	Validate() error
}

type Action[RequestType Request, ResponseType any] func(
	ctx context.Context, repository *core.Repository, r RequestType,
) (*ResponseType, error)
