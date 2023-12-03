CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS documents;

CREATE TABLE IF NOT EXISTS documents.t_user
(
    id        serial PRIMARY KEY,
    username  text NOT NULL,
    google_id text
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_google_id ON documents.t_user (google_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON documents.t_user (username);

CREATE TABLE IF NOT EXISTS documents.t_document
(
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        text   NOT NULL,
    owner       bigint NOT NULL REFERENCES documents.t_user (id),
    version     bigint NOT NULL  DEFAULT 1,
    description text
);

-- Start of the definition described by the session token management module.
-- One most likely should not change this definition.
CREATE TABLE IF NOT EXISTS sessions
(
    token  text PRIMARY KEY,
    data   bytea       NOT NULL,
    expiry timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions (expiry);
-- End;

DROP TYPE IF EXISTS link_status;
CREATE TYPE link_status AS ENUM ('enabled', 'disabled');

CREATE TABLE documents.t_link
(
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    document_id uuid NOT NULL REFERENCES documents.t_document (id) ON DELETE CASCADE,
    creation_dt timestamptz NOT NULL DEFAULT now(),
    expiry_dt timestamptz NOT NULL,
    status link_status NOT NULL DEFAULT 'enabled'
);
