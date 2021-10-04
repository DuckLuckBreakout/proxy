DROP DATABASE IF EXISTS proxy;
CREATE DATABASE proxy
    WITH OWNER postgres
    ENCODING 'utf8';
\connect proxy;

DROP TABLE IF EXISTS requests;
CREATE TABLE IF NOT EXISTS requests (
    id          SERIAL NOT NULL PRIMARY KEY,
    method      TEXT NOT NULL,
    scheme      TEXT NOT NULL,
    host        TEXT NOT NULL,
    path        TEXT NOT NULL,
    headers     JSONB NOT NULL,
    body        TEXT NOT NULL,
    params      JSONB NOT NULL
);