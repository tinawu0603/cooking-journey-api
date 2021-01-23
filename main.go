// main.go
package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/ziutek/mymysql/godrv"
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
	// Manually set DB_USERNAME and DB_PASSWORD environment variable
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")
	db_tcphost := "127.0.0.1"
	db_port := "3306"
	db_name := "cooking_journey"
	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_username, db_password, db_tcphost, db_port, db_name)

	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	handleRequests()
}
