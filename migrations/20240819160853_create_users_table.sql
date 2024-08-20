-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    login varchar(100) NOT NULL UNIQUE,
    password varchar(255) NOT NULL,
    created_at timestamp NOT NULL
);
CREATE INDEX login_idx ON users (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
