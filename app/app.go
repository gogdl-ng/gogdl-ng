package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/gorilla/mux"
)

func Run() {
	if err := gdrive.New(); err != nil {
		log.Fatalf("failed to initialize Google Drive service: %v", err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router = router.PathPrefix("/api/v1").Subrouter()

	router.HandleFunc("/jobs", CreateDownloadJob()).Methods("POST")

	go download.Run()

	log.Fatal(http.ListenAndServe(":3200", router))
}

func CreateDownloadJob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var job download.Job

		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		folder, err := gdrive.GetFolderById(job.DriveId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = download.RegisterNewJob(folder); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
