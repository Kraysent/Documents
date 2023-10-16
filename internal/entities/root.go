package entities

import (
	"fmt"
)

type Document struct {
	Username   string         `json:"username"`
	Type       string         `json:"document_type"`
	Attributes map[string]any `json:"attributes"`
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
