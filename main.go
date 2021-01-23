// main.go
package main

import (
	"fmt"
	"net/http"

	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var Recipes []Recipe
var db *sql.DB
var err error

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/recipes", getAllRecipesService).Methods("GET")
	myRouter.HandleFunc("/recipes", createNewRecipeService).Methods("POST")
	myRouter.HandleFunc("/recipes/{id}", deleteRecipe).Methods("DELETE")
	myRouter.HandleFunc("/recipes/{id}", getRecipeService).Methods("GET")
	myRouter.HandleFunc("/recipes/{id}", updateRecipeService).Methods("PUT")
	// log.Fatal(http.ListenAndServe(":10000", myRouter))
	http.ListenAndServe(":8000", myRouter)
}

func main() {
	db, err = sql.Open("mysql", "root:ocm5oShSbewYDeS6WcA2iqf7WnmYyx@tcp(127.0.0.1:3306)/cooking_journey?parseTime=true")
  if err != nil {
    panic(err.Error())
	}
	defer db.Close()
	handleRequests()
}
