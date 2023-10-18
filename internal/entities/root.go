package entities

import (
	"fmt"
)

type Document struct {
	ID         []byte         `json:"id" db:"id"`
	Username   string         `json:"username" db:"username"`
	Type       string         `json:"document_type" db:"document_type"`
	Attributes map[string]any `json:"attributes" db:"attributes"`
}

func (d *Document) Validate() error {
	if d.Username == "" {
		return fmt.Errorf("empty username")
	}

	if d.Type == "" {
		return fmt.Errorf("empty type")
	}

	return nil
}
