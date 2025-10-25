CREATE TABLE IF NOT EXISTS link_country_view_counter (
  id uuid,
  link_id uuid NOT NULL,
  country_code VARCHAR(2) NOT NULL,
  total_views INTEGER DEFAULT 0,
  unique_visitors INTEGER DEFAULT 0,
  created_on TIMESTAMPTZ DEFAULT NOW(),
  FOREIGN KEY (link_id) REFERENCES short_link(id) ON DELETE CASCADE,
  PRIMARY KEY (link_id, country_code)
);
