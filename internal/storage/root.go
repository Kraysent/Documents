package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type Storage interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	QuerySq(ctx context.Context, sqlizer sq.Sqlizer) (pgx.Rows, error)
	QueryRowSq(ctx context.Context, sqlizer sq.Sqlizer) (pgx.Row, error)
	ExecSq(ctx context.Context, sqlizer sq.Sqlizer) (int64, error)
}
