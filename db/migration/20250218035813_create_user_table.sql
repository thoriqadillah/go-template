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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
