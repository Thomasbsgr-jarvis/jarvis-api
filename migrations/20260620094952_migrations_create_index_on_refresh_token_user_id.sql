-- +goose Up
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);

-- +goose Down
DROP INDEX idx_refresh_tokens_user_id;
