create table users (
  id varchar(32) primary key,
  login varchar(32) unique,
  password varchar(64),
  name varchar(32),
  created_at timestamp,
  updated_at timestamp
);
