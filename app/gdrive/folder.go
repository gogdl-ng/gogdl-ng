package gdrive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

const FOLDER_MIMETYPE = "application/vnd.google-apps.folder"

func EnsureIsFolder(drive *drive.Service, folderId string) error {
	item, err := drive.Files.Get(folderId).Do()

	if err != nil {
		return err
	}

	if item.MimeType != FOLDER_MIMETYPE {
		return fmt.Errorf("item with id '%v' is not a folder", folderId)
	}

	return nil
}
