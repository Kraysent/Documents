package schema

import (
	"documents/internal/library/validation"
	"documents/internal/storage/links"
)

type GetLinksRequest struct {
	DocumentID string `schema:"document_id"`
	Status     string `schema:"status"`
	Page       uint64 `schema:"page"`
	PageSize   uint64 `schema:"page_size"`
}

func (r GetLinksRequest) Validate() error {
	rules := validation.NewRules(
		validation.In(
			r.Status, []string{string(links.StatusEnabled), string(links.StatusDisabled), ""},
		),
		validation.IsBetween(r.PageSize, 0, 100),
		validation.IsUUID(r.DocumentID),
	)

	return rules.Validate()
}

type Link struct {
	ID           string `json:"id"`
	DocumentID   string `json:"document_id"`
	CreationDate string `json:"creation_date"`
	ExpiryDate   string `json:"expiry_date"`
	Status       string `json:"status"`
}

type GetLinksResponse struct {
	Links []Link `json:"links"`
}
