package documents

import (
	"context"
	"database/sql"
	"fmt"

	libstorage "documents/internal/library/storage"
	"documents/internal/storage"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// DocumentStorage is a thin wrapper around documents table.
type DocumentStorage interface {
	AddDocument(context.Context, AddDocumentRequest) (*AddDocumentResult, error)
	RemoveDocuments(context.Context, map[string]any) (int64, error)
	GetDocuments(ctx context.Context, request GetDocumentsRequest) (*GetDocumentsResult, error)
	GetDocument(ctx context.Context, request GetDocumentRequest) (*GetDocumentResult, error)
}

type documentStorageImpl struct {
	storage storage.Storage
}

func NewDocumentStorage(store storage.Storage) *documentStorageImpl {
	return &documentStorageImpl{storage: store}
}

func (s *documentStorageImpl) AddDocument(
	ctx context.Context, request AddDocumentRequest,
) (result *AddDocumentResult, err error) {
	rows, err := s.storage.QuerySq(ctx,
		sq.Insert(TableName).
			Columns(ColumnName, ColumnOwner, ColumnDescription).
			Values(request.Name, request.Owner, request.Description).
			PlaceholderFormat(sq.Dollar).
			Suffix(fmt.Sprintf("RETURNING %s", ColumnID)))
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
		err = rows.Err()
	}()

	var documentID uuid.UUID
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		if err := rows.Scan(&documentID); err != nil {
			return nil, err
		}
	}

	return &AddDocumentResult{ID: documentID}, nil
}

func (s *documentStorageImpl) RemoveDocuments(ctx context.Context, fields map[string]any) (int64, error) {
	return s.storage.ExecSq(ctx, sq.Delete(TableName).
		Where(libstorage.SqAnd(fields)).
		PlaceholderFormat(sq.Dollar))
}

func (s *documentStorageImpl) GetDocument(
	ctx context.Context, request GetDocumentRequest,
) (*GetDocumentResult, error) {
	res, err := s.GetDocuments(ctx, GetDocumentsRequest{
		Fields: map[string]any{ColumnID: request.ID},
	})
	if err != nil {
		return nil, err
	}

	if len(res.Documents) == 0 {
		return nil, sql.ErrNoRows
	}

	if n := len(res.Documents); n != 1 {
		return nil, fmt.Errorf("unable to collect row, found %d rows", n)
	}

	return &GetDocumentResult{Document: res.Documents[0]}, nil
}

func (s *documentStorageImpl) GetDocuments(
	ctx context.Context, request GetDocumentsRequest,
) (result *GetDocumentsResult, err error) {
	rows, err := s.storage.QuerySq(ctx,
		sq.Select(ColumnID, ColumnName, ColumnOwner, ColumnVersion, ColumnDescription).
			From(TableName).
			Where(libstorage.SqAnd(request.Fields)).
			PlaceholderFormat(sq.Dollar))
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
		err = rows.Err()
	}()

	result = &GetDocumentsResult{}

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		doc, err := pgx.RowToStructByName[Document](rows)
		if err != nil {
			return nil, err
		}

		result.Documents = append(result.Documents, doc)
	}

	return result, nil
}
