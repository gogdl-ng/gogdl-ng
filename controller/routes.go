package controller

import (
	"net/http"

	controller "github.com/LegendaryB/gogdl-ng/controller/task"
	"github.com/LegendaryB/gogdl-ng/model/task"
	"github.com/gorilla/mux"
)

func AddRoutes(router *mux.Router, taskStore task.Store) {
	router.Use(jsonMiddleware)

	router.HandleFunc("/api/v1/tasks", controller.GetAllTasks(taskStore)).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{id}", controller.GetTask(taskStore)).Methods("GET")
	router.HandleFunc("/api/v1/tasks", controller.CreateTask(taskStore)).Methods("POST")
	router.HandleFunc("/api/v1/tasks/{id}", controller.DeleteTask(taskStore)).Methods("GET")
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
