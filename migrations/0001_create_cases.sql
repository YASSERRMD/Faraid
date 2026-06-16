-- +goose Up
-- +goose StatementBegin
CREATE TABLE cases (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    input      JSONB NOT NULL,
    result     JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cases;
-- +goose StatementEnd
