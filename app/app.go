package app

import (
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/controllers"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/middlewares"
	"github.com/LegendaryB/gogdl-ng/app/models/task"
	"github.com/LegendaryB/gogdl-ng/app/persistence"
	"github.com/gorilla/mux"
)

func Run() {
	db, err := persistence.NewDatabase()

	if err != nil {
		log.Fatalf("Failed to create database context: %v", err)
	}

	if err := task.NewRepository(db); err != nil {
		log.Fatalf("Failed to create task repository: %v", err)
	}

	if err := gdrive.New(); err != nil {
		log.Fatalf("Failed to create drive service: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	router.Use(middlewares.JSONMiddleware)
	controllers.AddRoutes(router)

	log.Fatal(http.ListenAndServe(":3200", router))
}
