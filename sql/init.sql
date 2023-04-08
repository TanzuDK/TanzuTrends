CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tweets (
    id BIGINT PRIMARY KEY,
    time TEXT,
    username TEXT,
    text TEXT,
    hashtags TEXT
);