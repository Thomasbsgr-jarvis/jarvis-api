-- +goose Up
CREATE TABLE eggo_files (
  id BIGSERIAL PRIMARY KEY,
  complaint_id UUID NOT NULL REFERENCES eggo_complaints (id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  hash TEXT NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (complaint_id, hash)
);

-- +goose Down
DROP TABLE eggo_files;
