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
		service.logger.Infof("starting to download file (name: %s, id: %s)", driveFile.Name, driveFile.Id)

		fp := filepath.Join(folderPath, driveFile.Name)
		file, err := service.getLocalFile(fp, driveFile.Size)

		if err != nil {
			service.logger.Errorf("failed to acquire local file: %v", err)
			return err
		}

		defer file.Close()

		request := service.drive.Files.Get(driveFile.Id).
			SupportsAllDrives(true).
			SupportsTeamDrives(true).
			AcknowledgeAbuse(service.conf.AcknowledgeAbuseFlag)

		fi, err := file.Stat()

		if err != nil {
			service.logger.Errorf("failed to stat() file. %v", err)
			return err
		}

		if fi.Size() == driveFile.Size {
			service.logger.Infof("file is already finished. skipping.")
			return nil
		}

		request.Header().Add("Range", fmt.Sprintf("bytes=%d-", fi.Size()))

		response, err := request.Download()

		if err != nil {
			service.logger.Errorf("failed to fetch file. %v", err)
			return err
		}

		_, err = io.Copy(file, response.Body)

		if err != nil {
			service.logger.Errorf("failed to write buffer to file. %v", err)
			return err
		}

		md5checksum, err := utils.GetMd5Checksum(fp)

		if err != nil {
			service.logger.Errorf("failed to calculate md5 checksum. %v", err)
			return err
		}

		if md5checksum != driveFile.Md5Checksum {
			err = errors.New("md5 checksum mismatch")
			service.logger.Errorf("the md5 checksum of the local file does not match checksum of the remote file. %v", err)
			return err
		}

		service.logger.Info("finished file download")
		return nil
	}, retry.Attempts(service.conf.RetryThreeshold))
}

func (service *DriveService) getLocalFile(path string, maxSize int64) (*os.File, error) {
	fi, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)

		if err != nil {
			service.logger.Errorf("failed to create file. %v", err)
			return nil, err
		}

		return f, nil
	}

	if fi.Size() > maxSize {
		service.logger.Warnf("local file size is greater than remote. file is probably corrupt. removing it..")
		os.Remove(path)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)

	if err != nil {
		service.logger.Errorf("failed to open file. %v", err)
		return nil, err
	}

	return f, nil
}
