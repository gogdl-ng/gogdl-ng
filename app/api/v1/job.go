package api

import (
	"encoding/json"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
)

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
