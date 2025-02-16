-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
                       uuid          UUID         DEFAULT uuid_generate_v4() PRIMARY KEY,
                       username      VARCHAR(50)  NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       balance       NUMERIC(20, 0) NOT NULL DEFAULT 0 CHECK (balance >= 0)
);
CREATE INDEX IF NOT EXISTS ind_users_uuid ON users (uuid);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users
-- +goose StatementEnd
