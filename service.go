package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func getAllRecipesService(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(getAllRecipesRepo())
}

func getRecipeService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	json.NewEncoder(w).Encode(getRecipeRepo(vars["name"]))
}

func createNewRecipeService(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	var recipe Recipe
	fmt.Println(reqBody)
	json.Unmarshal(reqBody, &recipe)
	fmt.Println(recipe)
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
	id, _ := strconv.Atoi(vars["id"])
	deleteRecipeRepo(id)
}
