package documents

import (
	"context"
	"fmt"

	libstorage "documents/internal/library/storage"
	"documents/internal/storage"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// DocumentStorage is a thin wrapper around documents table.
type DocumentStorage struct {
	storage storage.Storage
}

func NewDocumentStorage(store storage.Storage) *DocumentStorage {
	return &DocumentStorage{storage: store}
}

func (s *DocumentStorage) AddDocument(
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

func (s *DocumentStorage) RemoveDocuments(ctx context.Context, fields map[string]any) (int64, error) {
	return s.storage.ExecSq(ctx, sq.Delete(TableName).
		Where(libstorage.SqAnd(fields)).
		PlaceholderFormat(sq.Dollar))
}

func (s *DocumentStorage) GetDocuments(
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
