package router

import (
	"go-crud-new/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/newuser", middleware.CreateUser).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/users", middleware.GetAllUser).Methods("GET", "OPTIONS")

	return router
}
