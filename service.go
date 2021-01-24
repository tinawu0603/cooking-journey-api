package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func getAllRecipesService(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(getAllRecipesRepo())
}

func getRecipeService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(getRecipeRepo(vars["id"]))
}

func createNewRecipeService(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var recipe Recipe
	json.Unmarshal(reqBody, &recipe)
	recipe.DateCreated = time.Now()
	createNewRecipeRepo(recipe)
	json.NewEncoder(w).Encode(recipe)
}

func updateRecipeService(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var recipe Recipe
	json.Unmarshal(reqBody, &recipe)
	updateRecipeRepo(recipe)
	json.NewEncoder(w).Encode(recipe)
}

func deleteRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]
	deleteRecipeRepo(id)
}
