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
		file, err := service.getDestinationFile(fp, driveFile.Size)

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

func (service *DriveService) getDestinationFile(path string, maxSize int64) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	if stat.Size() > maxSize {
		service.logger.Warnf("Size of local file > size of remote file. File will be truncated because it is probably corrupted.")

		if err = truncate(file); err != nil {
			return nil, err
		}
	}

	return file, nil
}

func truncate(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return err
	}

	_, err := file.Seek(0, io.SeekStart)

	if err != nil {
		return err
	}

	return nil
}
