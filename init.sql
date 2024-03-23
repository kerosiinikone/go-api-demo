CREATE DATABASE godev;

\c godev

CREATE TABLE items
(
  id serial NOT NULL,
  name varchar(255)
);

ALTER TABLE items ADD CONSTRAINT id
  PRIMARY KEY (id);






