package controllers

import (
	"github.com/gorilla/mux"
)

func AddRoutes(router *mux.Router) {
	router.HandleFunc("/tasks", GetTasks()).Methods("GET")
	router.HandleFunc("/tasks/{id}", GetTask()).Methods("GET")
	router.HandleFunc("/tasks", CreateTask()).Methods("POST")
	router.HandleFunc("/tasks/{id}", DeleteTask()).Methods("GET")
}
