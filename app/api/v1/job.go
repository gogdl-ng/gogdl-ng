package api

import (
	"encoding/json"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
)

var logger = logging.NewLogger()

func CreateDownloadJob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var job download.Job

		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			logger.Errorf("failed to decode request json to object. %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		folder, err := gdrive.GetFolderById(job.DriveId)

		if err != nil {
			logger.Errorf("failed to get google drive folder. %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = download.RegisterNewJob(folder); err != nil {
			logger.Errorf("failed to register a new job. %w", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		logger.Infof("registered new job (driveId: %s, driveName: %s)", folder.Id, folder.Name)
	}
}
