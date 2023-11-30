package links

import (
	"context"
	"fmt"

	libstorage "documents/internal/library/storage"
	"documents/internal/storage"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// LinkStorage is a thin wrapper around users table
type LinkStorage interface {
	CreateLink(context.Context, CreateLinkRequest) (*CreateLinkResult, error)
	SetLinkStatus(context.Context, SetLinkStatusRequest) (*SetLinkStatusResult, error)
	GetLink(context.Context, GetLinkRequest) (*GetLinkResult, error)
	GetLinks(context.Context, GetLinksRequest) (*GetLinksResult, error)
}

type linkStorageImpl struct {
	storage storage.Storage
}

func NewLinkStorage(store storage.Storage) *linkStorageImpl {
	return &linkStorageImpl{storage: store}
}

func (s *linkStorageImpl) CreateLink(
	ctx context.Context, request CreateLinkRequest,
) (result *CreateLinkResult, err error) {
	row, err := s.storage.QueryRowSq(ctx,
		sq.Insert(TableName).
			Columns(ColumnDocumentID, ColumnExpiryDate).
			Values(request.DocumentID, request.ExpiryDate).
			Suffix(fmt.Sprintf("RETURNING %s", ColumnID)).
			PlaceholderFormat(sq.Dollar),
	)
	if err != nil {
		return nil, err
	}

	var linkID uuid.UUID
	if err := row.Scan(&linkID); err != nil {
		return nil, err
	}

	return &CreateLinkResult{ID: linkID}, nil
}

func (s *linkStorageImpl) SetLinkStatus(
	ctx context.Context, request SetLinkStatusRequest,
) (result *SetLinkStatusResult, err error) {
	n, err := s.storage.ExecSq(ctx,
		sq.Update(TableName).
			Set(ColumnStatus, request.Status).
			Where(sq.Eq{ColumnID: request.ID}).
			PlaceholderFormat(sq.Dollar),
	)
	if err != nil {
		return nil, err
	}
	if n != 0 {
		return nil, fmt.Errorf("error during setting status: updated %d rows", n)
	}

	return &SetLinkStatusResult{}, nil
}

func (s *linkStorageImpl) GetLink(
	ctx context.Context, request GetLinkRequest,
) (result *GetLinkResult, err error) {
	res, err := s.GetLinks(ctx, GetLinksRequest{
		Fields: map[string]any{ColumnID: request.ID}, PageSize: 1, OrderByField: ColumnID,
	})
	if err != nil {
		return nil, err
	}
	if n := len(res.Links); n != 1 {
		return nil, fmt.Errorf("unable to collect row, found %d rows", n)
	}

	return &GetLinkResult{Link: res.Links[0]}, nil
}

func (s *linkStorageImpl) GetLinks(
	ctx context.Context, request GetLinksRequest,
) (result *GetLinksResult, err error) {
	rows, err := s.storage.QuerySq(ctx,
		sq.Select(ColumnID, ColumnDocumentID, ColumnCreationDate, ColumnExpiryDate, ColumnStatus).
			From(TableName).
			Where(libstorage.SqAnd(request.Fields)).
			OrderBy(request.OrderByField).
			Offset(request.PageSize*request.PageNumber).
			Limit(request.PageSize),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
		err = rows.Err()
	}()

	result = &GetLinksResult{}

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		link, err := pgx.RowToStructByName[Link](rows)
		if err != nil {
			return nil, err
		}

		result.Links = append(result.Links, link)
	}

	return result, nil
}
