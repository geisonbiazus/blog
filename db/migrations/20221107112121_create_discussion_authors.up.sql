BEGIN;
CREATE TABLE IF NOT EXISTS discussion_authors(
   id uuid PRIMARY KEY,
   name VARCHAR,
   avatar_url VARCHAR
);
COMMIT;