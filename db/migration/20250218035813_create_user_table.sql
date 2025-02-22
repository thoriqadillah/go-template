-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id uuid default gen_random_uuid() not null primary key,
    email varchar(255) not null unique,
    password text,
    name varchar(255),
    source varchar(255),
    reset_token text,
    verified_at timestamp,
    created_at timestamp default now(),
    updated_at timestamp default now()
);

CREATE INDEX idx_users_id ON users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_id;
DROP TABLE users;
-- +goose StatementEnd
