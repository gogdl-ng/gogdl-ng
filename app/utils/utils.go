package utils

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

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
