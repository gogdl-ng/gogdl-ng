package gdrive

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/utils"
	"github.com/avast/retry-go"
	"google.golang.org/api/drive/v3"
)

func (service *DriveService) DownloadFile(folderPath string, driveFile *drive.File) error {
	return retry.Do(func() error {
		service.logger.Infof("File: %s", driveFile.Name)

		fp := filepath.Join(folderPath, driveFile.Name)
		file, err := service.getLocalFile(fp, driveFile.Size)

		if err != nil {
			service.logger.Errorf("Failed to acquire local file: %v", err)
			return err
		}

		defer file.Close()

		fi, err := file.Stat()

		if err != nil {
			service.logger.Errorf("Failed to stat() file. %v", err)
			return err
		}

		if fi.Size() == driveFile.Size {
			if err := service.compareChecksums(fp, driveFile.Md5Checksum); err != nil {
				return err
			}

			service.logger.Infof("Already completed. Skipping..")
			return nil
		}

		request := service.drive.Files.Get(driveFile.Id).
			SupportsAllDrives(true).
			SupportsTeamDrives(true)

		request.Header().Add("Range", fmt.Sprintf("bytes=%d-", fi.Size()))

		response, err := request.Download()

		if err != nil {
			service.logger.Errorf("Failed to fetch file. %v", err)
			return err
		}

		_, err = io.Copy(file, response.Body)

		if err != nil {
			service.logger.Errorf("Failed to write buffer to file. %v", err)
			return err
		}

		if err := service.compareChecksums(fp, driveFile.Md5Checksum); err != nil {
			return err
		}

		service.logger.Info("Finished file")
		return nil
	}, retry.Attempts(service.conf.RetryThreeshold))
}

func (service *DriveService) compareChecksums(localFilePath string, remoteFileChecksum string) error {
	localFileChecksum, err := utils.GetMd5Checksum(localFilePath)

	if err != nil {
		service.logger.Errorf("Failed to calculate md5 checksum. %v", err)
		return err
	}

	if localFileChecksum != remoteFileChecksum {
		err = errors.New("MD5 checksum mismatch")
		service.logger.Errorf("MD5 checksum of local file != MD5 checksum of remote file. %v", err)
		return err
	}

	service.logger.Infof("MD5 checksums are matching!")

	return nil
}

func (service *DriveService) getLocalFile(path string, maxSize int64) (*os.File, error) {
	fi, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)

		if err != nil {
			service.logger.Errorf("Failed to create file. %v", err)
			return nil, err
		}

		return file, nil
	}

	if fi.Size() > maxSize {
		service.logger.Warnf("Size of local file > size of remote file. File will be removed because it is probably corrupted.")
		os.Remove(path)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)

	if err != nil {
		service.logger.Errorf("Failed to open file at path: %s. %v", path, err)
		return nil, err
	}

	return file, nil
}
