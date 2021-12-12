package gdrive

import (
	"fmt"
)

const FOLDER_MIMETYPE = "application/vnd.google-apps.folder"

type DriveFolder struct {
	Id    string
	Name  string
	Files []DriveFile
}

func Folder(id string) (*DriveFolder, error) {
	fileCall := service.Files.Get(id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := fileCall.Do()

	if err != nil {
		return nil, err
	}

	if file.MimeType != FOLDER_MIMETYPE {
		return nil, fmt.Errorf("resource with id '%s' is not a folder", file.Id)
	}

	files, err := getFiles(file.Id)

	return &DriveFolder{
		Id:    file.Id,
		Name:  file.Name,
		Files: files,
	}, err
}
