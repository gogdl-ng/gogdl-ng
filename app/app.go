package app

import (
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/controllers"
	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/middlewares"
	"github.com/LegendaryB/gogdl-ng/app/persistence"
	"github.com/gorilla/mux"
)

func Run() {
	if err := persistence.NewDbContext(); err != nil {
		log.Fatalf("failed to initialize db context: %v", err)
	}

	if err := gdrive.New(); err != nil {
		log.Fatalf("failed to instantiate google drive service: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	router.Use(middlewares.JSONMiddleware)
	controllers.AddRoutes(router)

	go executeDownloadRoutine()

	log.Fatal(http.ListenAndServe(":3200", router))
}

func executeDownloadRoutine() {
	// todo: handle errors: skip?
	errch := make(chan error)
	go download.Start(errch)
	err := <-errch

	if err != nil {
		log.Fatalf("%v", err)
	}
}
