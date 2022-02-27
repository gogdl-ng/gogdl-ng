package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

func (s *DriveService) GetFiles(folder *drive.File) ([]*drive.File, error) {
	var children []*drive.File
	var nextPageToken string

	for {
		query := fmt.Sprintf("'%s' in parents and mimeType != '%s' and trashed=false", folder.Id, folderMimeType)

		fileList, err := s.requestFiles(query, nextPageToken)

		if err != nil {
			s.logger.Errorf("failed to execute Google Drive api request. %v", err)
			return nil, err
		}

		children = append(children, fileList.Files...)
		nextPageToken = fileList.NextPageToken

		if len(nextPageToken) == 0 {
			break
		}
	}

	return children, nil
}

func (s *DriveService) GetFolder(folderId string) (*drive.File, error) {
	file, err := s.requestFile(folderId)

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

func (s *DriveService) requestFile(id string) (*drive.File, error) {
	serviceGetCall := s.drive.Files.Get(id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := serviceGetCall.Do()

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *DriveService) requestFiles(query string, nextPageToken string) (*drive.FileList, error) {
	serviceListCall := s.drive.Files.List().
		OrderBy("name").
		PageSize(100).
		SupportsAllDrives(true).
		SupportsTeamDrives(true).
		IncludeItemsFromAllDrives(true).
		IncludeTeamDriveItems(true).
		Fields("nextPageToken, files(id, name, size, md5Checksum, mimeType, trashed)").
		Q(query)

	if len(nextPageToken) == 0 {
		serviceListCall.PageToken(nextPageToken)
	}

	fileList, err := serviceListCall.Do()

	if err != nil {
		return nil, err
	}

	return fileList, nil
}
