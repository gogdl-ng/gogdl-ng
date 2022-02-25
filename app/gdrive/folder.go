package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

const serviceQuery = "nextPageToken, files(id, name, size, md5Checksum, mimeType, trashed)"
const folderMimeType = "application/vnd.google-apps.folder"

func (service *DriveService) GetFilesFromFolder(folder *drive.File) ([]*drive.File, error) {
	var children []*drive.File
	var nextPageToken string

	for {
		query := fmt.Sprintf("'%s' in parents and mimeType != '%s' and trashed=false", folder.Id, folderMimeType)

		serviceListCall := service.drive.Files.List().
			OrderBy("name").
			PageSize(100).
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
			service.logger.Errorf("Failed to execute drive service call. %v", err)
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

func (service *DriveService) GetFolderById(folderId string) (*drive.File, error) {
	serviceGetCall := service.drive.Files.Get(folderId).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := serviceGetCall.Do()

	if err != nil {
		service.logger.Errorf("Failed to execute drive service call. %v", err)
		return nil, err
	}

	if file.MimeType != folderMimeType {
		err = fmt.Errorf("resource with id '%s' is not a folder", file.Id)
		service.logger.Error(err)

		return nil, err
	}

	return file, nil
}
