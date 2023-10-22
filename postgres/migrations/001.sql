CREATE SCHEMA IF NOT EXISTS documents;

CREATE TABLE IF NOT EXISTS documents.t_user (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    google_id TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_google_id ON documents.t_user (google_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON documents.t_user (username);

CREATE TABLE IF NOT EXISTS documents.t_documents (
    id BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES documents.t_user(id),
    document_type TEXT NOT NULL,
    attributes JSONB
);
