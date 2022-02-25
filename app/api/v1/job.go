package api

import (
	"encoding/json"
	"net/http"

	"github.com/LegendaryB/gogdl-ng/app/download"
	"github.com/LegendaryB/gogdl-ng/app/logging"
)

type JobController struct {
	logger     logging.Logger
	jobService *download.JobService
}

func NewJobController(logger logging.Logger, jobService *download.JobService) *JobController {
	return &JobController{logger: logger, jobService: jobService}
}

func (controller *JobController) CreateJob() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		CreateJobRequest := struct {
			DriveId string
		}{}

		if err := json.NewDecoder(r.Body).Decode(&CreateJobRequest); err != nil {
			controller.logger.Errorf("failed to decode request json to object. %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if len(CreateJobRequest.DriveId) == 0 {
			msg := ("property 'DriveId' has no value.")

			controller.logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
		}

		if err := controller.jobService.CreateJob(CreateJobRequest.DriveId); err != nil {
			controller.logger.Errorf("failed to register a new job. %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		controller.logger.Infof("registered new job (driveId: %s)", CreateJobRequest.DriveId)
	}
}
