CREATE TABLE users (
  id           SERIAL      PRIMARY KEY,
  name         VARCHAR(100) NOT NULL,
  email        VARCHAR(100) NOT NULL UNIQUE,
  password     VARCHAR      NOT NULL,
  role         VARCHAR(20)  NOT NULL DEFAULT 'writer',
  created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);