-- +goose Up
-- +goose StatementBegin

---- Create main keeper table which holds metadata, dates, type and encrypted content
CREATE TYPE entry_type AS ENUM ('credential', 'text', 'blob', 'card');

CREATE TABLE IF NOT EXISTS entries (
    id serial PRIMARY KEY,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    user_id integer,
    metadata varchar(255) NOT NULL,
    ent_type entry_type NOT NULL,
    payload bytea NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE entries;
-- +goose StatementEnd
