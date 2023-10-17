CREATE SCHEMA IF NOT EXISTS documents;

CREATE TABLE IF NOT EXISTS documents.t_documents (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    document_type TEXT NOT NULL,
    attributes JSONB
);

CREATE UNIQUE INDEX idx_username_document_type ON documents.t_documents (username, document_type);
