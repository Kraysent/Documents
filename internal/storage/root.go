package storage

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CollectionName = "documents"
)

type Config struct {
	Hosts        []string `yaml:"hosts"`
	Username     string   `yaml:"username"`
	PasswordEnv  string   `yaml:"password_env"`
	DBName       string   `yaml:"db_name"`
	CertFilePath string   `yaml:"cert_file_path"`
}

type Storage struct {
	client *mongo.Client
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
	url := fmt.Sprintf("mongodb://%s:%s@%s/%s?tls=true&tlsCaFile=%s",
		s.config.Username,
		password,
		strings.Join(s.config.Hosts, ","),
		s.config.DBName,
		s.config.CertFilePath)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	s.client = client

	return nil
}

func (s *Storage) Disconnect(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *Storage) AddDocument(ctx context.Context, attributes map[string]any) (string, error) {
	documentsCollection := s.client.Database(s.config.DBName).Collection(CollectionName)
	res, err := documentsCollection.InsertOne(ctx, attributes)
	if err != nil {
		return "", err
	}

	objectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("cannot cast inserted ID to ObjectID: %v", res.InsertedID)
	}
	return objectID.Hex(), nil
}

func (s *Storage) GetDocumentByID(ctx context.Context, id string) (map[string]any, error) {
	documentsCollection := s.client.Database(s.config.DBName).Collection(CollectionName)
	var result map[string]any

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return result, documentsCollection.FindOne(ctx, bson.D{{"_id", objectID}}).Decode(&result)
}
