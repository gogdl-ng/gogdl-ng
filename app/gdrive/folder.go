package gdrive

import (
	"fmt"
	"path/filepath"

	"google.golang.org/api/drive/v3"
)

const mimeTypeFolder = "application/vnd.google-apps.folder"

func (s *DriveService) GetFiles(folder *drive.File) ([]*DriveFile, error) {
	return s.getFiles(folder, "")
}

func (s *DriveService) getFiles(folder *drive.File, path string) ([]*DriveFile, error) {
	var driveFiles []*DriveFile
	var nextPageToken string

	for {
		query := fmt.Sprintf("'%s' in parents and trashed=false", folder.Id)

		fileList, err := s.requestFiles(query, nextPageToken)

		if err != nil {
			return nil, err
		}

		for _, driveFile := range fileList.Files {
			if !isDriveFolder(driveFile) {
				driveFiles = append(driveFiles, &DriveFile{
					Remote: driveFile,
					Path:   filepath.Join(path, driveFile.Name),
				})
				continue
			}

			files, err := s.getFiles(driveFile, filepath.Join(path, driveFile.Name))

			if err != nil {
				return nil, err
			}

			driveFiles = append(driveFiles, files...)
		}

		nextPageToken = fileList.NextPageToken

		if len(nextPageToken) == 0 {
			break
		}
	}

	return driveFiles, nil
}

func (s *DriveService) GetFolder(folderId string) (*drive.File, error) {
	driveFile, err := s.requestFile(folderId)

	if err != nil {
		return nil, err
	}

	if !isDriveFolder(driveFile) {
		err = fmt.Errorf("resource with id '%s' is not a folder", driveFile.Id)
		s.logger.Error(err)

		return nil, err
	}

	return driveFile, nil
}

func isDriveFolder(driveFile *drive.File) bool {
	return driveFile.MimeType == mimeTypeFolder
}

func (s *DriveService) requestFile(id string) (*drive.File, error) {
	serviceGetCall := s.drive.Files.Get(id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	file, err := serviceGetCall.Do()

	if err != nil {
		s.logger.Errorf("failed to execute Google Drive api request. %v", err)
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
		s.logger.Errorf("failed to execute Google Drive api request. %v", err)
		return nil, err
	}

	return fileList, nil
}
