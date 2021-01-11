package osext

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

func Md5Checksum(localPath string) (string, error) {
	f, err := os.Open(localPath)
	if err != nil {
		return "", errors.Wrapf(err, "can't open file (path: %s)", localPath)
	}

	h := md5.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", errors.Wrap(err, "copying buffer failed")
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

