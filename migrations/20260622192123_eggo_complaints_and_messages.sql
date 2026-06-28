-- +goose Up
CREATE TABLE eggo_complaints (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  folder_id TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE eggo_messages (
  id BIGSERIAL PRIMARY KEY,
  complaint_id UUID NOT NULL REFERENCES eggo_complaints (id) ON DELETE CASCADE,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE eggo_messages;

DROP TABLE eggo_complaints;
