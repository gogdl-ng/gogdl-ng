package download

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/gdrive"
	"github.com/LegendaryB/gogdl-ng/app/logging"
	"google.golang.org/api/drive/v3"
)

const (
	downloads  string = "downloads"
	completed  string = "completed"
	incomplete string = "incomplete"
)

type JobManager struct {
	logger logging.Logger
	drive  *gdrive.DriveService

	CompletedDirectoryPath  string
	IncompleteDirectoryPath string
}

type Job struct {
	Remote *drive.File

	DriveId string
	Path    string

	Files []*JobFile
}

type JobFile struct {
	Remote *drive.File
	Path   string
}

func NewJobManager(logger logging.Logger, drive *gdrive.DriveService) (*JobManager, error) {
	completedDirectoryPath, err := createDownloadsDirectory(completed)

	if err != nil {
		return nil, err
	}

	incompleteDirectoryPath, err := createDownloadsDirectory(incomplete)

	if err != nil {
		return nil, err
	}

	return &JobManager{
		logger:                  logger,
		drive:                   drive,
		CompletedDirectoryPath:  completedDirectoryPath,
		IncompleteDirectoryPath: incompleteDirectoryPath,
	}, nil
}

func (manager *JobManager) CreateJobPackage(job *Job) error {

	return nil
}

func (manager *JobManager) GetJobPackages() ([]*Job, error) {

	return nil, nil
}

func (manager *JobManager) FinishJobPackage() error {
	return nil
}

func createDownloadsDirectory(name string) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	path := filepath.Join(wd, downloads, name)

	if err := os.MkdirAll(path, 0644); err != nil {
		return "", err
	}

	return path, nil
}

func (manager *JobManager) getSubfolders(path string) ([]string, error) {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var subfolders []string

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		path := filepath.Join(path, item.Name())

		subfolders = append(subfolders, path)
	}

	return subfolders, nil
}
