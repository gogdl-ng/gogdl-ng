package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

const FOLDER_MIMETYPE = "application/vnd.google-apps.folder"

type DriveFolder struct {
	Id   string
	Name string
}

func Folder(id string) (*DriveFolder, error) {
	fileCall := Service.Files.Get(id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := fileCall.Do()

	if err != nil {
		return nil, err
	}

	err = ensureIsFolder(file)

	if err != nil {
		return nil, err
	}

	return &DriveFolder{
		Id:   file.Id,
		Name: file.Name,
	}, err
}

func ensureIsFolder(file *drive.File) error {
	if file.MimeType != FOLDER_MIMETYPE {
		return fmt.Errorf("resource with id '%s' is not a folder", file.Id)
	}

	return nil
}
