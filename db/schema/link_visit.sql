CREATE TABLE IF NOT EXISTS link_visit (
  id uuid PRIMARY KEY,
  link_id uuid NOT NULL,
  ip INET,
  user_agent TEXT,
  created_on TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (link_id) REFERENCES short_link(id) ON DELETE CASCADE
);
