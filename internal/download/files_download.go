package download

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"
)

type fileState struct {
	file      *driveFile
	localPath string

	size        int64
	md5Checksum string
}

func (fs *fileState) valid() (bool, error) {
	if fs.size < fs.file.Size {
		return false, errors.New("file is not complete")
	}

	if fs.md5Checksum != fs.file.Md5Checksum {
		return false, errors.New("file is corrupted")
	}
	return true, nil
}

func (fs *fileState) offset() int64 {
	if fs.size > fs.file.Size {
		return 0 // file is corrupted (too long)
	}

	if fs.size == fs.file.Size && fs.md5Checksum != fs.file.Md5Checksum {
		return 0 // file is corrupted
	}
	return fs.size
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

	state.md5Checksum, err = calculateMd5Checksum(localPath)
	if err != nil {
		return nil, errors.Wrap(err, "calculating MD5 checksum failed")
	}

	state.size = fi.Size()
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

func downloadFile(driveService *drive.Service, file *driveFile, outputDir string) error {
	bar := progressbar.DefaultBytes(
		file.Size,
		renderBarDescription(file),
	)

	return retry.Do(func() error {
		state, err := evaluateFileState(file, filepath.Join(outputDir, file.Path))
		if err != nil {
			return errors.Wrap(err, "state evaluation failed")
		}

		if ok, _ := state.valid(); ok {
			bar.Finish()
			return nil // file is complete and valid
		}
		bar.Set64(state.offset())

		err = downloadFileData(driveService, *state, bar)
		if err != nil {
			return errors.Wrap(err, "downloading file data failed")
		}

		state, err = evaluateFileState(file, filepath.Join(outputDir, file.Path))
		if err != nil {
			return errors.Wrap(err, "state evaluation failed")
		}

		if ok, err := state.valid(); !ok {
			return errors.Wrap(err, "file verification failed")
		}

		return nil
	}, retry.Attempts(3))
}

func downloadFileData(driveService *drive.Service, state fileState, bar *progressbar.ProgressBar) error {
	filesGetCall := driveService.Files.Get(state.file.Id)
	filesGetCall.Header().Add("Range", fmt.Sprintf("bytes=%d-", state.offset()))
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
	var flags = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	if state.offset() == 0 {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	f, err := os.OpenFile(state.localPath, flags, 0644)
	if err != nil {
		return errors.Wrap(err, "opening destination file failed")
	}
	defer f.Close()

	// Download file data
	_, err = io.Copy(io.MultiWriter(f, bar), resp.Body)
	if err != nil {
		return errors.Wrap(err, "copying buffers failed")
	}
	return nil
}

func renderBarDescription(file *driveFile) string {
	var description strings.Builder
	if len(file.Name) > 22 {
		description.WriteString(file.Name[:22])
		description.WriteString("...")
	} else {
		description.WriteString(file.Name)
	}

	description.WriteString(" (MD5: ")
	description.WriteString(file.Md5Checksum)
	description.WriteByte(')')
	return description.String()
}
