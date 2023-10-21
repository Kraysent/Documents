package core

import (
	"documents/internal/storage"
	"documents/internal/storage/documents"
	"documents/internal/storage/users"
)

type Repository struct {
	Storage  storage.Storage
	Storages struct {
		Documents *documents.DocumentStorage
		Users     *users.UserStorage
	}
	Config *Config
}

func NewRepository(config *Config) (*Repository, error) {
	repo := &Repository{
		Config: config,
	}

	store := storage.NewStorage(config.Storage)

	repo.Storage = store
	repo.Storages.Documents = documents.NewDocumentStorage(store)
	repo.Storages.Users = users.NewUserStorage(store)

	return repo, nil
}
