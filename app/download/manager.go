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

type JobManager struct {
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

func NewJobManager(logger logging.Logger, conf *config.Configuration, drive *gdrive.DriveService) (*JobManager, error) {
	completedDirectoryPath, err := createDownloadsDirectory(completed)

	if err != nil {
		return nil, err
	}

	incompleteDirectoryPath, err := createDownloadsDirectory(incomplete)

	if err != nil {
		return nil, err
	}

	service := &JobManager{
		logger:                  logger,
		drive:                   drive,
		CompletedDirectoryPath:  completedDirectoryPath,
		IncompleteDirectoryPath: incompleteDirectoryPath,
	}

	service.dispatcher = NewDispatcher(service, conf.Queue.MaxWorkers, conf.Queue.Size)

	return service, nil
}

func (jm *JobManager) Run() error {
	unfinishedJobs, err := jm.GetUnfinishedJobs()

	if err != nil {
		return err
	}

	jm.dispatcher.Start(context.Background())
	jm.dispatcher.Wait()

	// todo: what when unfinished jobs > queueSize??
	jm.dispatcher.AddJobs(unfinishedJobs)

	return nil
}

func (jm *JobManager) RunJob(job *Job) {
	files, err := jm.drive.GetFilesFromFolder(job.File)

	if err != nil {
		jm.logger.Errorf("failed to retrieve files of folder: '%s'. %v", job.Id, err)
		return
	}

	for _, driveFile := range files {
		path := jm.getFileTargetPath(job, driveFile)

		if err := jm.drive.DownloadFile(driveFile, path); err != nil {
			jm.logger.Errorf("failed to download file (name: %s, id: %s). %v", driveFile.Name, driveFile.Id, err)
		}
	}

	jm.FinishJob(job)
}

func (jm *JobManager) CreateJob(driveId string) error {
	folder, err := jm.drive.GetFolderById(driveId)

	if err != nil {
		return err
	}

	path, err := jm.createJobDirectory(folder)

	if err != nil {
		return err
	}

	job := &Job{
		Path: path,
		File: folder,
	}

	jm.dispatcher.AddJob(job)

	return nil
}

func (jm *JobManager) GetUnfinishedJobs() ([]*Job, error) {

	return nil, nil
}

func (jm *JobManager) FinishJob(job *Job) error {
	if err := jm.removeDriveIdFile(job.Path); err != nil {
		return err
	}

	if err := jm.MoveToCompletedDirectory(job); err != nil {
		return err
	}

	return nil
}
