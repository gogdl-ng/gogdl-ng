package api

import (
	"encoding/json"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/logging"
)

type JobController struct {
	logger     logging.Logger
	downloader *download.Downloader
}

func NewJobController(logger logging.Logger, downloader *download.Downloader) *JobController {
	return &JobController{logger: logger, downloader: downloader}
}

func (controller *JobController) CreateDownloadJob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var job download.Job

		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			controller.logger.Errorf("failed to decode request json to object. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := controller.downloader.RegisterNewJob(job.DriveId); err != nil {
			controller.logger.Errorf("failed to register a new job. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		controller.logger.Infof("registered new job (driveId: %s)", job.DriveId)
	}
}
