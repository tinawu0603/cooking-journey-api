// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/ziutek/mymysql/godrv"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var db *sql.DB
var err error

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/recipes", getAllRecipesService).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/api/recipes", createNewRecipeService).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/api/recipes/{id}", deleteRecipe).Methods("DELETE", "OPTIONS")
	myRouter.HandleFunc("/api/recipes/{id}", getRecipeService).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/api/recipes/{id}", updateRecipeService).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	var isDev bool
	if len(os.Args) > 1 {
		isDev = true
	} else {
		isDev = false
	}
	var dbURI string
	if isDev {
		fmt.Println("YOU ARE RUNNING IN DEV MODE!")
		// DEV MODE:
		// Running Google SQL Proxy to connect to database
		dbUsername := "db-user"
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := "127.0.0.1"
		dbPort := "3306"
		dbName := "cooking_journey"
		dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbPort, dbName)
	} else {
		// CLOUD MODE:
		// GAE instance connecting to Google SQL server
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
		socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
		if !isSet {
			socketDir = "/cloudsql"
		}
		dbUsername := "db-user"
		dbPassword := string(result.Payload.Data)
		dbName := "cooking_journey"
		dbConnectionName := "plated-mechanic-302602:us-east4:database-1"
		dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUsername, dbPassword, socketDir, dbConnectionName, dbName)
	}

	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	handleRequests()
}
