CREATE TABLE IF NOT EXISTS "user" (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    oauth_provider VARCHAR(50) NOT NULL,
    oauth_subject VARCHAR(255) NOT NULL,
    created_on TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_oauth_identity UNIQUE (oauth_provider, oauth_subject)
);

CREATE INDEX idx_user_email ON "user"(email);
CREATE INDEX idx_user_oauth_identity ON "user"(oauth_provider, oauth_subject);
