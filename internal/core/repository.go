package core

import (
	"documents/internal/storage"
)

type Repository struct {
	Storage struct {
		DocumentStorage *storage.Storage
	}
	Config *Config
}

func NewRepository(config *Config) (*Repository, error) {
	repo := &Repository{
		Config: config,
	}

	repo.Storage.DocumentStorage = storage.NewStorage(config.Storage)

	return repo, nil
}
