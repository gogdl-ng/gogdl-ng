package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/LegendaryB/gogdl-ng/app/logging"
)

var logger = logging.NewLogger()

func Move(source string, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		logger.Errorf("failed to create folder(s). %v", err)
		return err
	}

	items, err := ioutil.ReadDir(source)

	if err != nil {
		logger.Errorf("failed to read folder content. %v", err)
		return err
	}

	for _, item := range items {
		sourcefp := filepath.Join(source, item.Name())
		targetfp := filepath.Join(target, item.Name())

		if err := os.Rename(sourcefp, targetfp); err != nil {
			logger.Errorf("failed to move file. %v", err)
			return err
		}
	}

	if err = os.Remove(source); err != nil {
		logger.Errorf("failed to delete folder. %v", err)
		return err
	}

	return nil
}

func GetMd5Checksum(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		logger.Errorf("failed to open file. %v", err)
		return "", err
	}

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		logger.Errorf("failed to write buffer. %v", err)
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func Subfolders(path string) ([]fs.FileInfo, error) {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		logger.Errorf("failed to read folder content. %v", err)
		return nil, err
	}

	var subfolders []fs.FileInfo

	for _, item := range items {
		if !item.IsDir() {
			logger.Infof("%s is no folder, skipping.", item.Name())
			continue
		}

		subfolders = append(subfolders, item)
	}

	return subfolders, nil
}
