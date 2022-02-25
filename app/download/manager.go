package download

import (
	"github.com/LegendaryB/gogdl-ng/app/config"
	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

const (
	downloads       string = "downloads"
	completed       string = "completed"
	incomplete      string = "incomplete"
	driveIdFileName string = "driveId"
)

type JobService struct {
	logger     logging.Logger
	drive      *gdrive.DriveService
	dispatcher *Dispatcher

	CompletedDirectoryPath  string
	IncompleteDirectoryPath string
}

type Worker interface {
	RunJob(job *Job)
}

type Job struct {
	Path string
	*drive.File
}

func NewJobService(logger logging.Logger, conf *config.Configuration, drive *gdrive.DriveService) (*JobService, error) {
	completedDirectoryPath, err := createDownloadsDirectory(completed)

	if err != nil {
		return nil, err
	}

	incompleteDirectoryPath, err := createDownloadsDirectory(incomplete)

	if err != nil {
		return nil, err
	}

	service := &JobService{
		logger:                  logger,
		drive:                   drive,
		CompletedDirectoryPath:  completedDirectoryPath,
		IncompleteDirectoryPath: incompleteDirectoryPath,
	}

	service.dispatcher = NewDispatcher(service, conf.Jobs.MaxWorkers, conf.Jobs.QueueSize)

	return service, nil
}

func (service *JobService) Run() error {
	unfinishedJobs, err := service.GetUnfinishedJobs()

	if err != nil {
		return err
	}

	service.dispatcher.Start(context.Background())
	service.dispatcher.Wait()

	// todo: what when unfinished jobs > queueSize??
	service.dispatcher.AddJobs(unfinishedJobs)

	return nil
}

func (service *JobService) RunJob(job *Job) {
	files, err := service.drive.GetFilesFromFolder(job.File)

	if err != nil {
		service.logger.Error("") // todo: log
		return
	}

	for _, driveFile := range files {
		path := service.getFileTargetPath(job, driveFile)

		if err := service.drive.DownloadFile(driveFile, path); err != nil {
			service.logger.Errorf("Failed to download file (name: %s, id: %s). %v", driveFile.Name, driveFile.Id, err)
		}
	}

	service.FinishJob(job)
}

func (service *JobService) CreateJob(driveId string) error {
	folder, err := service.drive.GetFolderById(driveId)

	if err != nil {
		return err
	}

	path, err := service.createJobDirectory(folder)

	if err != nil {
		return err
	}

	job := &Job{
		Path: path,
		File: folder,
	}

	service.dispatcher.AddJob(job)

	return nil
}

func (service *JobService) GetUnfinishedJobs() ([]*Job, error) {

	return nil, nil
}

func (service *JobService) FinishJob(job *Job) error {
	if err := service.removeDriveIdFile(job.Path); err != nil {
		return err
	}

	if err := service.MoveToCompletedDirectory(job); err != nil {
		return err
	}

	return nil
}
