package storage

import (
	"context"
	"fmt"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

const (
	DocumentsTableName          = "documents.t_documents"
	DocumentsColumnID           = "id"
	DocumentsColumnUsername     = "username"
	DocumentsColumnDocumentType = "document_type"
	DocumentsColumnAttributes   = "attributes"
)

type Config struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Username    string `yaml:"username"`
	PasswordEnv string `yaml:"password_env"`
	DBName      string `yaml:"db_name"`
}

type Storage struct {
	db     *pgx.Conn
	config Config
}

func NewStorage(config Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	password, ok := os.LookupEnv(s.config.PasswordEnv)
	if !ok {
		return fmt.Errorf("no DB password specified on %s env", s.config.PasswordEnv)
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.config.Host,
		s.config.Port,
		s.config.Username,
		password,
		s.config.DBName,
	)

	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(ctx); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Disconnect(ctx context.Context) error {
	return s.db.Close(ctx)
}

func (s *Storage) QuerySq(ctx context.Context, query sq.Sqlizer) (pgx.Rows, error) {
	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *Storage) ExecSq(ctx context.Context, query sq.Sqlizer) (int64, error) {
	q, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	result, err := s.db.Exec(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}
