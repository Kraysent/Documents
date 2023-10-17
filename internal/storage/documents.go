package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"documents/internal/entities"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) AddDocument(
	ctx context.Context, username, documentType string, attributes map[string]any,
) (rowsAffected int64, err error) {
	attrBytes, err := json.Marshal(attributes)
	if err != nil {
		return 0, err
	}

	query := sq.Insert(DocumentsTableName).
		Columns(DocumentsColumnUsername, DocumentsColumnDocumentType, DocumentsColumnAttributes).
		Values(username, documentType, attrBytes).
		PlaceholderFormat(sq.Dollar).
		Suffix(fmt.Sprintf("RETURNING %s", DocumentsColumnID))

	rows, err := s.QuerySq(ctx, query)
	if err != nil {
		return 0, err
	}
	defer func() {
		rows.Close()
		err = rows.Err()
	}()

	var documentID int64
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return 0, err
		}

		if err := rows.Scan(&documentID); err != nil {
			return 0, err
		}
	}

	return documentID, nil
}

func (s *Storage) RemoveDocument(ctx context.Context, fields map[string]any) (int64, error) {
	selector := sq.And{}

	for key, value := range fields {
		selector = append(selector, sq.Eq{key: value})
	}

	query := sq.Delete(DocumentsTableName).
		Where(selector).
		PlaceholderFormat(sq.Dollar)

	return s.ExecSq(ctx, query)
}

func (s *Storage) GetDocuments(ctx context.Context, request GetDocumentsRequest) (*GetDocumentsResult, error) {
	selector := sq.And{}

	for key, value := range request.Fields {
		selector = append(selector, sq.Eq{key: value})
	}

	query := sq.Select(DocumentsColumnID, DocumentsColumnUsername, DocumentsColumnDocumentType, DocumentsColumnAttributes).
		From(DocumentsTableName).
		Where(selector).
		PlaceholderFormat(sq.Dollar)

	rows, err := s.QuerySq(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		rows.Close()
		err = rows.Err()
	}()

	var result GetDocumentsResult

	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, err
		}

		doc, err := pgx.RowToStructByName[entities.Document](rows)
		if err != nil {
			return nil, err
		}

		result.Documents = append(result.Documents, doc)
	}

	return &result, nil
}
