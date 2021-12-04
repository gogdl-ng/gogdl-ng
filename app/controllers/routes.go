package controllers

import (
	"github.com/LegendaryB/gogdl-ng/app/persistence"

	"github.com/gorilla/mux"
)

func AddRoutes(router *mux.Router, dbContext *persistence.DbContext) {
	router.HandleFunc("/tasks", GetAllTasks(dbContext.Tasks)).Methods("GET")
	router.HandleFunc("/tasks/{id}", GetTask(dbContext.Tasks)).Methods("GET")
	router.HandleFunc("/tasks", CreateTask(dbContext.Tasks)).Methods("POST")
	router.HandleFunc("/tasks/{id}", DeleteTask(dbContext.Tasks)).Methods("GET")
}
