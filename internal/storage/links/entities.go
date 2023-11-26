package links

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusEnabled  Status = "enabled"
	StatusDisabled Status = "disabled"
)

const (
	TableName          = "documents.t_link"
	ColumnID           = "id"
	ColumnDocumentID   = "document_id"
	ColumnCreationDate = "creation_dt"
	ColumnExpiryDate   = "expiry_dt"
	ColumnStatus       = "status"
)

type Link struct {
	ID           uuid.UUID `db:"id"`
	DocumentID   uuid.UUID `db:"document_id"`
	CreationDate time.Time `db:"creation_dt"`
	ExpiryDate   time.Time `db:"expiry_date"`
	Status       Status    `db:"status"`
}

type CreateLinkRequest struct {
	DocumentID uuid.UUID
	ExpiryDate time.Time
}

type CreateLinkResult struct {
	ID uuid.UUID
}

type SetLinkStatusRequest struct {
	ID     uuid.UUID
	Status Status
}

type SetLinkStatusResult struct{}

type GetLinkRequest struct {
	ID uuid.UUID
}

type GetLinkResult struct {
	Link Link
}

type GetLinksRequest struct {
	Fields       map[string]any
	PageSize     uint64
	PageNumber   uint64
	OrderByField string
}

type GetLinksResult struct {
	Links []Link
}
