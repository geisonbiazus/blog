BEGIN;
CREATE TABLE IF NOT EXISTS users(
   id uuid PRIMARY KEY,
   provider_user_id VARCHAR UNIQUE,
   email VARCHAR UNIQUE,
   name VARCHAR,
   avatar_url VARCHAR
);
COMMIT;