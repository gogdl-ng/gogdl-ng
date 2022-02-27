package download

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

func (jm *JobManager) createJobDirectory(driveFolder *drive.File) (string, error) {
	path := filepath.Join(jm.IncompleteDirectoryPath, driveFolder.Name)

	if err := os.MkdirAll(path, 0644); err != nil {
		return "", err
	}

	if err := jm.createDriveIdFile(path, driveFolder.Id); err != nil {
		return "", err
	}

	return path, nil
}

func (jm *JobManager) getFileTargetPath(job *Job, driveFile *drive.File) string {
	return filepath.Join(job.Path, driveFile.Name)
}

func createDownloadsDirectory(folderName string) (string, error) {
	wd, err := os.Getwd()

	if err != nil {
		return "", err
	}

	path := filepath.Join(wd, downloads, folderName)

	if err := os.MkdirAll(path, 0644); err != nil {
		return "", err
	}

	return path, nil
}

func (jm *JobManager) moveToCompletedDirectory(job *Job) error {
	targetPath := filepath.Join(jm.CompletedDirectoryPath, filepath.Base(job.Path))

	if err := os.MkdirAll(targetPath, 0644); err != nil {
		return err
	}

	items, err := ioutil.ReadDir(job.Path)

	if err != nil {
		return err
	}

	for _, item := range items {
		sourcefp := filepath.Join(job.Path, item.Name())
		targetfp := filepath.Join(targetPath, item.Name())

		if err := os.Rename(sourcefp, targetfp); err != nil {
			return err
		}
	}

	if err = os.Remove(job.Path); err != nil {
		return err
	}

	return nil
}

func (jm *JobManager) getSubfolders(path string) ([]string, error) {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var names []string

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		names = append(names, filepath.Join(path, item.Name()))
	}

	return names, nil
}
