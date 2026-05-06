CREATE TABLE users (
    id varchar(32) PRIMARY KEY,
    login varchar(32) UNIQUE,
    password varchar(64),
    name varchar(32),
    created_at timestamp, updated_at timestamp
);
