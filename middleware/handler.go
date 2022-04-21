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
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id, omitempty"`
	Message string `json:"message, omitempty"`
}

// use godotenv, os, sql
func createConnection() *sql.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRESS_URL"))

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfuly connected!")

	return db
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
		ID:      insertID,
		Message: "user created Successfuly",
	}

	json.NewEncoder(w).Encode(res)

}

//  int64 is a return value's type
func insertUser(user models.User) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO users(firstname, lastname, email) VALUES ($1, $2, $3) RETURNING userid`

	var id int64

	err := db.QueryRow(sqlStatement, user.Firstname, user.Lastname, user.Email).Scan(&id)

	if err != nil {
		log.Fatal("Unable excute the query %v", err)
	}

	fmt.Println("Insert a single record %v", id)

	return id
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// to use mux.Vars you can get params
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("couldnt convert string into integer %v", err)
	}

	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("couldnt get user %v", err)
	}

	json.NewEncoder(w).Encode(user)

}

func getUser(id int64) (models.User, error) {
	db := createConnection()

	defer db.Close()

	var user models.User

	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows retuned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable scan row %v", err)

	}

	return user, err
}

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	user, err := getAllUsers()

	if err != nil {
		log.Fatalf("couldnt get user %v", err)
	}

	json.NewEncoder(w).Encode(user)
}

func getAllUsers() ([]models.User, error) {
	db := createConnection()

	defer db.Close()

	var users []models.User

	sqlStatement := `SELECT * FROM users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("couldnt get users %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email)

		if err != nil {
			log.Fatalf("couldnt scan row %v", err)
		}

		users = append(users, user)
	}

	return users, err
}
