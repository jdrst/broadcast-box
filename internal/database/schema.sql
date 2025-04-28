CREATE TABLE users (
  id   TEXT       PRIMARY KEY,
  name TEXT       NOT NULL UNIQUE,
  password  BLOB  NOT NULL,
  streamKey BLOB  NOT NULL UNIQUE
);
