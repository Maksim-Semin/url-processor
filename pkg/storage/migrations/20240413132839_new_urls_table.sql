-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS URLs (
    id SERIAL PRIMARY KEY,
    shortCode VARCHAR(10) NOT NULL,
    URL VARCHAR(255) NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS URLs;
-- +goose StatementEnd
