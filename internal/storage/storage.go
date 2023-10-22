package storage

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	PasswordEnv     string `yaml:"password_env"`
	DBName          string `yaml:"db_name"`
	SSLMode         string `yaml:"ssl_mode"`
	SSLRootCertPath string `yaml:"ssl_root_cert_path"`
}

type storageImpl struct {
	DB     *pgxpool.Pool
	Config Config
}

func NewStorage(config Config) *storageImpl {
	return &storageImpl{
		Config: config,
	}
}

func (s *storageImpl) Connect(ctx context.Context) error {
	password, ok := os.LookupEnv(s.Config.PasswordEnv)
	if !ok {
		return fmt.Errorf("no DB password specified on %s env", s.Config.PasswordEnv)
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.Config.Host,
		s.Config.Port,
		s.Config.Username,
		password,
		s.Config.DBName,
		s.Config.SSLMode,
	)

	if s.Config.SSLMode == "require" {
		dsn += fmt.Sprintf(" sslrootcert=%s", s.Config.SSLRootCertPath)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	s.DB = pool

	return nil
}

func (s *storageImpl) Disconnect(ctx context.Context) error {
	s.DB.Close()

	return nil
}

func (s *storageImpl) QuerySq(ctx context.Context, query sq.Sqlizer) (pgx.Rows, error) {
	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *storageImpl) QueryRowSq(ctx context.Context, query sq.Sqlizer) (pgx.Row, error) {
	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := s.DB.QueryRow(ctx, q, args...)

	return row, nil
}

func (s *storageImpl) ExecSq(ctx context.Context, query sq.Sqlizer) (int64, error) {
	q, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	result, err := s.DB.Exec(ctx, q, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func RandomHex(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}
