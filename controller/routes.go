package controller

import (
	"net/http"

	controller "github.com/LegendaryB/gogdl-ng/controller/task"
	"github.com/LegendaryB/gogdl-ng/interfaces"
	"github.com/gorilla/mux"
)

func AddRoutes(router *mux.Router, repository interfaces.ITaskRepository) {
	router.Use(jsonMiddleware)

	router.HandleFunc("/api/v1/tasks", controller.GetAllTasks(repository)).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{id}", controller.GetTask(repository)).Methods("GET")
	router.HandleFunc("/api/v1/tasks", controller.CreateTask(repository)).Methods("POST")
	router.HandleFunc("/api/v1/tasks/{id}", controller.DeleteTask(repository)).Methods("GET")
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
