package gdrive

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/LegendaryB/gogdl-ng/app/utils"
	"github.com/avast/retry-go"
	"google.golang.org/api/drive/v3"
)

func (ds *DriveService) DownloadFile(driveFile *drive.File, path string) error {
	return retry.Do(func() error {
		ds.logger.Infof("file: %s", driveFile.Name)

		file, err := ds.getDestinationFile(path, driveFile.Size)

		if err != nil {
			ds.logger.Errorf("failed to get handle of destination file: %v", err)
			return err
		}

		defer file.Close()

		fi, err := file.Stat()

		if err != nil {
			ds.logger.Errorf("failed to stat() file. %v", err)
			return err
		}

		if fi.Size() == driveFile.Size {
			if err := ds.compareChecksums(path, driveFile.Md5Checksum); err != nil {
				return err
			}

			ds.logger.Infof("Already completed. Skipping..")
			return nil
		}

		response, err := ds.requestFileContent(driveFile, fi.Size())

		if err != nil {
			return err
		}

		_, err = io.Copy(file, response.Body)

		if err != nil {
			ds.logger.Errorf("Failed to write buffer to file. %v", err)
			return err
		}

		if err := ds.compareChecksums(path, driveFile.Md5Checksum); err != nil {
			return err
		}

		ds.logger.Info("Finished file")
		return nil
	}, retry.Attempts(ds.conf.Download.RetryThreeshold))
}

func (ds *DriveService) requestFileContent(driveFile *drive.File, rangeStart int64) (*http.Response, error) {
	request := ds.drive.Files.Get(driveFile.Id).
		SupportsAllDrives(true).
		SupportsTeamDrives(true)

	request.Header().Add("Range", fmt.Sprintf("bytes=%d-", rangeStart))

	response, err := request.Download()

	if err != nil {
		ds.logger.Errorf("failed to fetch content of file. %v", err)
	}

	return response, err
}

func (ds *DriveService) compareChecksums(localFilePath string, remoteFileChecksum string) error {
	localFileChecksum, err := utils.GetMd5Checksum(localFilePath)

	if err != nil {
		ds.logger.Errorf("failed to calculate md5 checksum. %v", err)
		return err
	}

	if localFileChecksum != remoteFileChecksum {
		err = errors.New("MD5 checksum mismatch")
		ds.logger.Errorf("MD5 checksum of local file != MD5 checksum of remote file. %v", err)
		return err
	}

	ds.logger.Infof("MD5 checksums are matching!")

	return nil
}

func (ds *DriveService) getDestinationFile(path string, maxSize int64) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()

	if err != nil {
		return nil, err
	}

	if stat.Size() > maxSize {
		ds.logger.Warnf("size of local file > size of remote file. file will be truncated because it is probably corrupted.")

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
