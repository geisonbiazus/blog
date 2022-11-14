BEGIN;
CREATE TABLE IF NOT EXISTS discussion_comments(
   id uuid PRIMARY KEY,
   subject_id VARCHAR,
   author_id uuid REFERENCES discussion_authors(id),
   markdown text,
   html text,
   created_at timestamp without time zone
);
COMMIT;