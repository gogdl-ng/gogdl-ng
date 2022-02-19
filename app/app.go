package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/api/v1"
	"github.com/LegendaryB/gogdl-ng/app/config"
	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"github.com/gorilla/mux"
)

func Run() {
	logger := logging.NewLogger()

	if err := env.InitializeEnvironment(); err != nil {
		logger.Fatalf("failed to initialize environment. %w", err)
	}

	if err := config.LoadConfiguration(); err != nil {
		logger.Fatalf("failed to retrieve application configuration. %w", err)
	}

	if err := gdrive.New(); err != nil {
		logger.Fatalf("failed to initialize Google Drive service. %w", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	registerRoutes(router)

	go listenAndServe(router, config.Loaded.Application.ListenPort)

	if err := download.Run(); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/jobs", api.CreateDownloadJob()).Methods("POST")
}

func listenAndServe(router *mux.Router, listenPort int) {
	addr := fmt.Sprintf(":%d", listenPort)

	log.Fatal(http.ListenAndServe(addr, router))
}
