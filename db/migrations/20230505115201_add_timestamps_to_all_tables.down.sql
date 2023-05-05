BEGIN;
ALTER TABLE auth_users DROP COLUMN created_at;
ALTER TABLE auth_users DROP COLUMN updated_at;

ALTER TABLE discussion_authors DROP COLUMN created_at;
ALTER TABLE discussion_authors DROP COLUMN updated_at;

ALTER TABLE discussion_comments DROP COLUMN updated_at;
END;