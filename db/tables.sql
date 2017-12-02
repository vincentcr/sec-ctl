BEGIN;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users(
  id uuid PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);

DROP TABLE IF EXISTS devices CASCADE;
CREATE TABLE sites(
  id uuid PRIMARY KEY,
  owner_id uuid REFERENCES users(id) ON DELETE RESTRICT,
  state_shadow JSONB
);


DROP TABLE IF EXISTS auth_tokens CASCADE;
CREATE TABLE auth_tokens(
  token uuid PRIMARY KEY,
  rec_id uuid NOT NULL,
  expires_at TIMESTAMP
);


DROP TABLE IF EXISTS events CASCADE;
CREATE TABLE events(
  site_id uuid PRIMARY KEY,
  level TEXT NOT NULL,
  time TIMESTAMP NOT NULL,
  data JSONB
);

COMMIT;
