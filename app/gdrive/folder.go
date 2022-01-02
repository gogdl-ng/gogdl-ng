package gdrive

import (
	"errors"
	"fmt"
	"path/filepath"
)

const serviceQuery = "nextPageToken, files(id, name, md5Checksum, mimeType, trashed)"
const serviceMaxPageSize = 100

type DriveFolder struct {
	Id    string
	Name  string
	Files []*DriveFile
}

func GetFilesByFolderId(folderId string) (*DriveFolder, error) {
	folder, err := getFolderById(folderId)

	if err != nil {
		return nil, err
	}

	files, err := getFilesByPath(folderId, "")

	if err != nil {
		return nil, err
	}

	folder.Files = files

	return folder, nil
}

func getFilesByPath(folderId string, path string) ([]*DriveFile, error) {
	var files []*DriveFile
	var nextPageToken string

	for {
		query := fmt.Sprintf("'%s' in parents and trashed=false", folderId)

		serviceListCall := service.Files.List().
			PageSize(serviceMaxPageSize).
			OrderBy("name").
			SupportsAllDrives(true).
			SupportsTeamDrives(true).
			IncludeItemsFromAllDrives(true).
			IncludeTeamDriveItems(true).
			Fields(serviceQuery).
			Q(query)

		if len(nextPageToken) == 0 {
			serviceListCall.PageToken(nextPageToken)
		}

		list, err := serviceListCall.Do()

		if err != nil {
			return nil, errors.New("failed to execute service list call")
		}

		for _, file := range list.Files {
			if file.MimeType == Folder {
				childFolderFiles, err := getFilesByPath(file.Id, filepath.Join(path, file.Name))

				if err != nil {
					return nil, fmt.Errorf("failed to retrieve child folder (id: %s) files", file.Id)
				}

				files = append(files, childFolderFiles...)
				continue
			}

			files = append(files, &DriveFile{
				Id:       file.Id,
				Name:     file.Name,
				Path:     filepath.Join(path, file.Name),
				Checksum: file.Md5Checksum,
			})
		}

		nextPageToken = list.NextPageToken

		if len(nextPageToken) == 0 {
			break
		}
	}

	return files, nil
}

func getFolderById(folderId string) (*DriveFolder, error) {
	serviceGetCall := service.Files.Get(folderId).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := serviceGetCall.Do()

	if err != nil {
		return nil, err
	}

	if file.MimeType != Folder {
		return nil, fmt.Errorf("resource with id '%s' is not a folder", file.Id)
	}

	return &DriveFolder{
		Id:   file.Id,
		Name: file.Name,
	}, err
}
