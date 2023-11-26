package core

import (
	"time"

	"documents/internal/storage"
	"documents/internal/storage/documents"
	"documents/internal/storage/links"
	"documents/internal/storage/users"
	"github.com/alexedwards/scs/v2"
)

type Repository struct {
	Storage  storage.Storage
	Storages struct {
		Documents *documents.DocumentStorage
		Users     *users.UserStorage
		Links     *links.LinkStorage
	}
	SessionManager *scs.SessionManager
	Config         *Config
}

func NewRepository(config *Config) (*Repository, error) {
	repo := &Repository{
		Config: config,
	}

	repo.Storage = storage.NewStorage(config.Storage)
	repo.Storages.Documents = documents.NewDocumentStorage(repo.Storage)
	repo.Storages.Users = users.NewUserStorage(repo.Storage)
	repo.Storages.Links = links.NewLinkStorage(repo.Storage)

	repo.SessionManager = scs.New()
	repo.SessionManager.Lifetime = 24 * time.Hour

	return repo, nil
}
