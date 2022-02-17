package app

import (
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/api/v1"
	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"github.com/gorilla/mux"
)

func Run() {
	logger := logging.NewLogger()

	if err := gdrive.New(); err != nil {
		logger.Fatalf("failed to initialize Google Drive service: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	registerRoutes(router)

	go listenAndServe(router)

	download.Run()
}

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/jobs", api.CreateDownloadJob()).Methods("POST")
}

func listenAndServe(router *mux.Router) {
	log.Fatal(http.ListenAndServe(":3200", router))
}
