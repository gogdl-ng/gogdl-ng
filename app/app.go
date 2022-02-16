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
		log.Fatalf("failed to initialize Google Drive service: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	router.Use(middlewares.JSONMiddleware)
	controllers.AddRoutes(router)

	go download.Run()

	log.Fatal(http.ListenAndServe(":3200", router))
}
