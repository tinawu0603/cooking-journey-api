package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func getAllRecipesRepo() []Recipe {
	var recipes []Recipe
	// Get all recipes
	recipesResult, err := db.Query("SELECT * from recipes")
	if err != nil {
		panic(err.Error())
	}
	defer recipesResult.Close()

	// Get all ingredients in each recipe
	for recipesResult.Next() {
		var rec Recipe
		err := recipesResult.Scan(&rec.Id, &rec.Name, &rec.DateCreated)
		if err != nil {
			panic(err.Error())
		}
		rec.Ingredients = getIngredientsForRecipeRepo(rec.Id)
		recipes = append(recipes, rec)
	}

	return recipes
}

func getRecipeRepo(name string) Recipe {
	var recipe Recipe
	recipeResult, err := db.Query("SELECT * from recipes WHERE name = ?", name)
	if err != nil {
		panic(err.Error())
	}
	defer recipeResult.Close()
	for recipeResult.Next() {
		err := recipeResult.Scan(&recipe.Id, &recipe.Name, &recipe.DateCreated)
		if err != nil {
			panic(err.Error())
		}
		recipe.Ingredients = getIngredientsForRecipeRepo(recipe.Id)
	}
	return recipe
}

func createNewRecipeRepo(recipe Recipe) {
	recipesStmt, err := db.Prepare("INSERT INTO recipes(name, dateCreated) VALUES(?, ?)")
	if err != nil {
		panic(err.Error())
	}
	if recipeAlreadyExist(recipe.Name) {
		fmt.Println("Recipe already existed, please choose a different name.")
	} else {
		// Recipe does not exist yet, create it
		_, err = recipesStmt.Exec(recipe.Name, recipe.DateCreated)
		if err != nil {
			panic(err.Error())
		}
		// Create each Ingredient in this recipe
		// First, have to get the recipe's generated ID from recipes table
		recipeResult, err := db.Query("SELECT id from recipes WHERE name = ?", recipe.Name)
		if err != nil {
			panic(err.Error())
		}
		defer recipeResult.Close()
		for recipeResult.Next() {
			err := recipeResult.Scan(&recipe.Id)
			if err != nil {
				panic(err.Error())
			}
		}
		for _, ing := range recipe.Ingredients {
			ingStmt, err := db.Prepare("INSERT INTO ingredients(name, quantity, unit, recipeId) VALUES(?, ?, ?, ?)")
			if err != nil {
				panic(err.Error())
			}
			_, err = ingStmt.Exec(ing.Name, ing.Quantity, ing.Unit, recipe.Id)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func updateRecipeRepo(recipe Recipe) {
	// TODO: Implement UPSERT behavior in MySQL
	// Source: https://www.techbeamers.com/mysql-upsert/
	recipesStmt, err := db.Prepare("UPDATE recipes SET name = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = recipesStmt.Exec(recipe.Name, recipe.Id)
	if err != nil {
		panic(err.Error())
	}
	// Update each Ingredient in this recipe
	for _, ing := range recipe.Ingredients {
		ingStmt, err := db.Prepare("UPDATE ingredients SET name = ?, quantity = ?, unit = ? WHERE id = ?")
		if err != nil {
			panic(err.Error())
		}
		_, err = ingStmt.Exec(ing.Name, ing.Quantity, ing.Unit, ing.Id)
		if err != nil {
			panic(err.Error())
		}
	}
}

func deleteRecipeRepo(recipeId int) {
	recipesStmt, err := db.Prepare("DELETE FROM recipes WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = recipesStmt.Exec(recipeId)
	if err != nil {
		panic(err.Error())
	}
	ingStmt, err := db.Prepare("DELETE FROM ingredients WHERE recipeId = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = ingStmt.Exec(recipeId)
	if err != nil {
		panic(err.Error())
	}
}

/*
 *	HELPERS
 */

// Return a list of Ingredients for the given Recipe
func getIngredientsForRecipeRepo(recipeId int) []Ingredient {
	var ingredients []Ingredient
	ingredientsResult, err := db.Query("SELECT id, name, quantity, unit from ingredients WHERE recipeId = ?", recipeId)
	if err != nil {
		panic(err.Error())
	}
	defer ingredientsResult.Close()
	for ingredientsResult.Next() {
		var ing Ingredient
		err := ingredientsResult.Scan(&ing.Id, &ing.Name, &ing.Quantity, &ing.Unit)
		if err != nil {
			panic(err.Error())
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients
}

// Check if the given Recipe name already exists
func recipeAlreadyExist(name string) bool {
	row := db.QueryRow("SELECT name FROM recipes WHERE name = ?", name)
	switch err := row.Scan(&name); err {
	case sql.ErrNoRows:
		return false
	case nil: // Recipe exists
		return true
	default:
		panic(err.Error())
	}
}
