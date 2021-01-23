CREATE DATABASE cooking_journey;

CREATE TABLE recipes (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(36) NOT NULL,
  dateCreated date NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE ingredients (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(36) NOT NULL,
  quantity smallint NOT NULL,
  unit varchar(16) NOT NULL,
  recipeId int NOT NULL REFERENCES recipes(id),
  PRIMARY KEY (id)
);