package download

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func getMd5Checksum(path string) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", fmt.Errorf("can't open file - path: %s", path)
	}

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
