package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

const serviceQuery = "nextPageToken, files(id, name, size, md5Checksum, mimeType, trashed)"
const folderMimeType = "application/vnd.google-apps.folder"

func (s *DriveService) GetFiles(folder *drive.File) ([]*drive.File, error) {
	var children []*drive.File
	var nextPageToken string

	for {
		query := fmt.Sprintf("'%s' in parents and mimeType != '%s' and trashed=false", folder.Id, folderMimeType)

		serviceListCall := s.drive.Files.List().
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
			s.logger.Errorf("failed to execute Google Drive api request. %v", err)
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

func (s *DriveService) GetFolder(folderId string) (*drive.File, error) {
	serviceGetCall := s.drive.Files.Get(folderId).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := serviceGetCall.Do()

	if err != nil {
		s.logger.Errorf("failed to execute Google Drive api request. %v", err)
		return nil, err
	}

	if file.MimeType != folderMimeType {
		err = fmt.Errorf("resource with id '%s' is not a folder", file.Id)
		s.logger.Error(err)

		return nil, err
	}

	return file, nil
}
