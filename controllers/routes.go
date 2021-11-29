package controllers

import (
	"net/http"

	controllers "github.com/LegendaryB/gogdl-ng/controllers/task"
	"github.com/LegendaryB/gogdl-ng/interfaces"
	"github.com/gorilla/mux"
)

func AddRoutes(router *mux.Router, repository interfaces.ITaskRepository) {
	router.Use(jsonMiddleware)

	router.HandleFunc("/api/v1/tasks", controllers.GetAllTasks(repository)).Methods("GET")
	router.HandleFunc("/api/v1/tasks/{id}", controllers.GetTask(repository)).Methods("GET")
	router.HandleFunc("/api/v1/tasks", controllers.CreateTask(repository)).Methods("POST")
	router.HandleFunc("/api/v1/tasks/{id}", controllers.DeleteTask(repository)).Methods("GET")
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
