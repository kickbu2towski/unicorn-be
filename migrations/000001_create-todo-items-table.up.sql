CREATE TABLE IF NOT EXISTS todo_items (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  name text NOT NULL,
  state text NOT NULL,
  priority text NOT NULL,
  Tags text[],
  closed_at timestamp(0) with time zone
);