package download

import (
	"github.com/avast/retry-go"
	"github.com/pkg/errors"
)

type fileState struct {
	file *driveFile

	localPath   string
	size        int64
	md5Checksum string

	err error
}

func (fs *fileState) verify() (bool, error) {
	// TODO 	Check if local file exists in output directory
	// TODO 	Check MD5 remote vs local file
	// TODO 	Check if size(local file) < size(remote local)
	return false, nil // TODO
}

func evaluateFileState(file *driveFile, localDir string) fileState {
	return fileState{} // TODO
}

func processFile(state fileState) error {
	if valid, _ := state.verify(); valid {
		return nil // the file is complete and valid
	}
	return retry.Do(func() error {
		err := downloadFile(state)
		if err != nil {
			return errors.Wrap(err, "downloading file failed")
		}

		if valid, err := state.verify(); !valid {
			return errors.Wrap(err, "file verification failed")
		}
		return nil
	}, retry.Attempts(3))
}

func downloadFile(state fileState) error {
	return nil // TODO
}
