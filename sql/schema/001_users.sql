-- +goose Up
CREATE TABLE users (
    id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    CONSTRAINT users_pk PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE users;
