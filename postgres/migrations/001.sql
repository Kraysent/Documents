CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SCHEMA IF NOT EXISTS documents;

CREATE TABLE IF NOT EXISTS documents.t_user
(
    id        SERIAL PRIMARY KEY,
    username  TEXT NOT NULL,
    google_id TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_google_id ON documents.t_user (google_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_username ON documents.t_user (username);

-- CREATE TYPE attribute AS
-- (
--     id    text,
--     name  text,
--     value text
-- );

CREATE TABLE IF NOT EXISTS documents.t_document
(
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        text   NOT NULL,
    owner       bigint NOT NULL REFERENCES documents.t_user (id),
    version     bigint NOT NULL DEFAULT 1,
    description text
--     attributes  attribute[]
);
