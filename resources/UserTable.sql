CREATE TABLE IF NOT EXISTS users(
  id       serial primary key,
  name     varchar(256) NOT NULL,
  password varchar(256) NOT NULL
);

INSERT INTO users(name, password)
  VALUES ('Test1', '123');

SELECT * FROM users;
