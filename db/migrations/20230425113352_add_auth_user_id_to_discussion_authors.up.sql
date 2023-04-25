BEGIN;
ALTER TABLE discussion_authors ADD COLUMN IF NOT EXISTS auth_user_id UUID;
CREATE INDEX IF NOT EXISTS discussion_authors_auth_user_id_idx ON discussion_authors(auth_user_id);
COMMIT;
