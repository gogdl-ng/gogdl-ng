package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/api/v1"
	"github.com/LegendaryB/gogdl-ng/app/config"
	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/env"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"github.com/gorilla/mux"
)

func Run() {
	conf, err := config.NewConfigurationFromFile()

	if err != nil {
		log.Fatalf("failed to retrieve app configuration. %v", err)
	}

	logger, err := logging.NewLogger1(conf.Application.LogFilePath)

	if err != nil {
		log.Fatalf("failed to initialize logger. %s", err)
	}

	if err := env.NewEnvironment(logger); err != nil {
		logger.Fatalf("failed to initialize environment. %v", err)
	}

	downloader, err := download.NewDownloader(&conf.Transfer, logger)

	if err != nil {
		logger.Fatalf("failed to initialize Downloader service. %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	jobController := api.NewJobController(logger, downloader)

	router.HandleFunc("/jobs", jobController.CreateDownloadJob()).Methods("POST")

	go listenAndServe(router, conf.Application.ListenPort)

	if err := downloader.Run(); err != nil {
		logger.Fatal(err)
	}
}

func listenAndServe(router *mux.Router, listenPort int) {
	addr := fmt.Sprintf(":%d", listenPort)

	log.Fatal(http.ListenAndServe(addr, router))
}
