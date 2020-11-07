package download

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"
)

type fileState struct {
	file      *driveFile
	localPath string

	offset      int64
	md5Checksum string
}

func (fs *fileState) verify() (bool, error) {
	if fs.offset != fs.file.Size {
		return false, errors.New("file is not complete")
	}

	if fs.md5Checksum != fs.file.Md5Checksum {
		return false, errors.New("file is corrupted")
	}
	return true, nil
}

func evaluateFileState(file *driveFile, localPath string) (*fileState, error) {
	state := &fileState{
		file:      file,
		localPath: localPath,
	}

	fi, err := os.Stat(localPath)
	if os.IsNotExist(err) {
		return state, nil // file hasn't been downloaded yet
	}
	if err != nil {
		return nil, errors.Wrapf(err, "stat file failed (path: %s)", state.localPath)
	}

	state.offset = fi.Size()
	state.md5Checksum, err = calculateMd5Checksum(localPath)
	if err != nil {
		return nil, errors.Wrap(err, "calculating MD5 checksum failed")
	}

	if state.md5Checksum != state.file.Md5Checksum {
		state.offset = 0
	}
	return state, nil
}

func calculateMd5Checksum(localPath string) (string, error) {
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

func processFile(driveService *drive.Service, file *driveFile, outputDir string) error {
	bar := progressbar.DefaultBytes(
		file.Size,
		fmt.Sprintf("%s (MD5: %s)", file.Name, file.Md5Checksum),
	)

	return retry.Do(func() error {
		state, err := evaluateFileState(file, filepath.Join(outputDir, file.Path))
		if err != nil {
			return errors.Wrap(err, "state evaluation failed")
		}

		if valid, _ := state.verify(); valid {
			bar.Finish()
			return nil // file is complete and valid
		}

		err = downloadFile(driveService, *state, bar)
		if err != nil {
			return errors.Wrap(err, "downloading file failed")
		}

		if valid, err := state.verify(); !valid {
			return errors.Wrap(err, "file verification failed")
		}

		bar.Finish()
		return nil
	}, retry.Attempts(3))
}

func downloadFile(driveService *drive.Service, state fileState, bar *progressbar.ProgressBar) error {
	filesGetCall := driveService.Files.Get(state.file.Id)
	filesGetCall.Header().Add("Range", fmt.Sprintf("bytes=%d-", state.offset))
	resp, err := filesGetCall.Download()
	if err != nil {
		return errors.Wrap(err, "opening file stream failed")
	}

	// Create base directories
	basePath := filepath.Dir(state.localPath)
	err = os.MkdirAll(basePath, 0755)
	if err != nil {
		return errors.Wrapf(err, "creating base directories failed (path: %s)", basePath)
	}

	// Open destination file
	f, err := os.OpenFile(state.localPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return errors.Wrap(err, "opening destination file failed")
	}
	defer f.Close()

	// Download file data
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return errors.Wrap(err, "downloading file data failed")
	}
	return nil
}
