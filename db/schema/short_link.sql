CREATE TABLE short_link (
  id uuid PRIMARY KEY,
  original_route TEXT NOT NULL,
  client_id uuid NULL,
  code VARCHAR(10) UNIQUE NOT NULL,
  created_on TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE
);
