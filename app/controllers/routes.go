package controllers

import (
	"github.com/LegendaryB/gogdl-ng/app/models/task"
	"github.com/gorilla/mux"
)

func AddRoutes(router *mux.Router, repository *task.Repository) {
	router.HandleFunc("/tasks", GetAllTasks(repository)).Methods("GET")
	router.HandleFunc("/tasks/{id}", GetTask(repository)).Methods("GET")
	router.HandleFunc("/tasks", CreateTask(repository)).Methods("POST")
	router.HandleFunc("/tasks/{id}", DeleteTask(repository)).Methods("GET")
}
