package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

const serviceQuery = "nextPageToken, files(id, name, size, md5Checksum, mimeType, trashed)"
const folderMimeType = "application/vnd.google-apps.folder"
const maxPageSize = 100

type DriveFolder struct {
	Id   string
	Name string
}

func (service *DriveService) GetFilesFromFolder(folderId string) ([]*drive.File, error) {
	var children []*drive.File
	var nextPageToken string

	for {
		query := fmt.Sprintf("'%s' in parents and mimeType != '%s' and trashed=false", folderId, folderMimeType)

		serviceListCall := service.drive.Files.List().
			PageSize(maxPageSize).
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
			service.logger.Errorf("failed to execute drive service call. %v", err)
			return nil, err
		}

		children = append(children, list.Files...)

		nextPageToken = list.NextPageToken

		if len(nextPageToken) == 0 {
			break
		}
	}

	return children, nil
}

func (service *DriveService) GetFolderById(folderId string) (*DriveFolder, error) {
	serviceGetCall := service.drive.Files.Get(folderId).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := serviceGetCall.Do()

	if err != nil {
		service.logger.Errorf("failed to execute drive service call. %v", err)
		return nil, err
	}

	if file.MimeType != folderMimeType {
		err = fmt.Errorf("resource with id '%s' is not a folder", file.Id)
		service.logger.Error(err)

		return nil, err
	}

	return &DriveFolder{
		Id:   file.Id,
		Name: file.Name,
	}, err
}
