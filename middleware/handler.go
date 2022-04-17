package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-crud-new/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
)

type responce struct {
	ID int64 `json:"id, omitempty"`
	Message string `json:"message, omitempty"`
}

// use godotenv, os, sql
func createConnection() *sql.DB {
	err := godotenv.Load()

	if err != nil{
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgress", os.Getenv("POSTGRESS_URL"))

	if err != nil{
		panic(err)
	}

	fmt.Println("Successfuly connected!")

	return db
}

func CreateUser (w http.ResponseWriter, r *http.Request){
	w.Header().Set("Context-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatal("Unable to decode the request body. %v", err)
	}

	insertID := insertUser(user)

	res := response{
		ID: insertID,
		Message: "user created Successfuly",
	}


	json.NewEncoder(w).Encode(res)

}

//  int64 is a return value's type
func insertUser(user models.User) ini64{
	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO users(firstname, lastname, email) VALUES ($1, $2, $3) RETURNING userid`

	var id int64

	err := db.QueryLow(user.Firstname, user.Lastname, user.Email).Scan(&id)

	if err != nil{
		log.Fatal("Unable excute the query %v", err)
	}

	fmt.Println("Insert a single record %v", id)

	return id
}

