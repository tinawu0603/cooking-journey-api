// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/ziutek/mymysql/godrv"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
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
	// Get database password from Google secret manager
	secretName := "projects/793806054736/secrets/db-password/versions/1"
	// Create the client
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		fmt.Println("failed to create secretmanager client")
		panic(err.Error())
	}
	// Build the request
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}
	// Call the API
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		fmt.Println("failed to access secret version")
		panic(err.Error())
	}
	// Manually set DB_PASSWORD environment variable
	// LOCAL DEV
	// db_password := os.Getenv("DB_PASSWORD")
	// db_tcphost := "127.0.0.1"

	// GOOGLE CLOUD
	dbPassword := string(result.Payload.Data)
	//dbHost := "35.245.37.12"
	dbUsername := "db-user"
	//dbPort := "3306"
	dbName := "cooking_journey"
	dbConnectionName := "plated-mechanic-302602:us-east4:database-1"
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}
	var dbURI string
	// dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUsername, dbPassword, socketDir, dbConnectionName, dbName)

	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	handleRequests()
}
