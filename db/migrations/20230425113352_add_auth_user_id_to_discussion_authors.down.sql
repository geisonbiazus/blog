BEGIN;
ALTER TABLE discussion_authors DROP COLUMN IF EXISTS auth_user_id;
COMMIT;
