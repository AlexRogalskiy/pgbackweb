-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS executions (
  id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,

  status TEXT NOT NULL CHECK (
    status IN ('running', 'success', 'failed', 'deleted')
  ) DEFAULT 'running',
  message TEXT NOT NULL,
  path TEXT,

  started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ,
  finished_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

CREATE TRIGGER executions_change_updated_at
BEFORE UPDATE ON executions FOR EACH ROW EXECUTE FUNCTION change_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS executions;
-- +goose StatementEnd