package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Move(source string, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	items, err := ioutil.ReadDir(source)

	if err != nil {
		return err
	}

	for _, item := range items {
		sourcefp := filepath.Join(source, item.Name())
		targetfp := filepath.Join(target, item.Name())

		if err := os.Rename(sourcefp, targetfp); err != nil {
			return err
		}
	}

	if err = os.Remove(source); err != nil {
		return err
	}

	return nil
}

func GetMd5Checksum(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func Subfolders(path string) ([]fs.FileInfo, error) {
	items, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var subfolders []fs.FileInfo

	for _, item := range items {
		if !item.IsDir() {
			continue
		}

		subfolders = append(subfolders, item)
	}

	return subfolders, nil
}
