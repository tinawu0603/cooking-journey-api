package main

import (
	"time"
)

/*
Table: recipes
+-------------+-------------+------+-----+---------+----------------+
| Field       | Type        | Null | Key | Default | Extra          |
+-------------+-------------+------+-----+---------+----------------+
| id          | int         | NO   | PRI | NULL    | auto_increment |
| name        | varchar(36) | NO   |     | NULL    |                |
| dateCreated | date        | NO   |     | NULL    |                |
+-------------+-------------+------+-----+---------+----------------+
*/

/*
Table: ingredients
+----------+-------------+------+-----+---------+----------------+
| Field    | Type        | Null | Key | Default | Extra          |
+----------+-------------+------+-----+---------+----------------+
| id       | int         | NO   | PRI | NULL    | auto_increment |
| name     | varchar(36) | NO   |     | NULL    |                |
| quantity | smallint    | YES  |     | NULL    |                |
| unit     | varchar(16) | YES  |     | NULL    |                |
| recipeId | int         | YES  | MUL | NULL    |                |
+----------+-------------+------+-----+---------+----------------+
*/

// Ingredient - struct for all ingredients
type Ingredient struct {
	Id				int `json:Id`
	Name			string `json:"Name"`
	Quantity	int `json:"Quantity"`
	Unit			string `json:"Unit"`			
}

// Recipe - struct for all recipes
type Recipe struct {
	Id					int `json:"Id"`
	Name				string `json:"Name"`
	Ingredients []Ingredient `json:"Ingredients"`
	DateCreated time.Time
}
