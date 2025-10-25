CREATE TABLE IF NOT EXISTS link_analytics (
  id uuid PRIMARY KEY,
  link_id uuid NOT NULL,
  total_views INTEGER DEFAULT 0,
  unique_visitors INTEGER DEFAULT 0,
  created_on TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (link_id) REFERENCES short_link(id) ON DELETE CASCADE
);
