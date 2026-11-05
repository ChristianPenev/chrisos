-- +goose Up
CREATE TABLE users (
    id         uuid        PRIMARY KEY,
    email      text        NOT NULL CONSTRAINT users_email_key UNIQUE,
    first_name text        NOT NULL,
    last_name  text        NOT NULL,
    created_at timestamptz NOT NULL,
    CONSTRAINT users_names_not_empty CHECK (first_name <> '' AND last_name <> '')
);

-- +goose Down
DROP TABLE users;
